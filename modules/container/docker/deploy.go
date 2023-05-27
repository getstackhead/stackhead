package container_docker

import (
	"encoding/json"
	"fmt"
	"path"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	diff_docker_compose "github.com/saitho/diff-docker-compose/lib"
	xfs "github.com/saitho/golang-extended-fs/v2"
	logger "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	container_docker_definitions "github.com/getstackhead/stackhead/modules/container/docker/definitions"
	docker_compose "github.com/getstackhead/stackhead/modules/container/docker/docker-compose"
	docker_system "github.com/getstackhead/stackhead/modules/container/docker/system"
	"github.com/getstackhead/stackhead/system"
)

type Data struct {
	Context     system.ContextStruct
	DockerPaths container_docker_definitions.DockerPaths
}

type RegExSettings struct {
	Pattern string
	Repl    string
}

var userNameCache = map[string]int{}

func resolveUserNameWithCache(userName string) (int, error) {
	value, ok := userNameCache[userName]
	if ok {
		return value, nil
	}

	var resolvedUserId int
	var err error
	if reflect.TypeOf(userName).String() == "string" { // if userId is a string, resolve it to uid
		resolvedUserId, err = system.ResolveRemoteUserIntoUid(userName)
		if err != nil {
			// check if string is a number
			if userAsInt, err2 := strconv.Atoi(userName); err2 == nil {
				userNameCache[userName] = userAsInt
				return userAsInt, nil
			}
		}
	} else {
		resolvedUserId, err = strconv.Atoi(userName)
	}
	userNameCache[userName] = resolvedUserId
	return resolvedUserId, err
}

func (m Module) Deploy(modulesSettings interface{}) error {
	dockerPaths := container_docker_definitions.GetDockerPaths()
	// Build src folder list
	srcFolderList := container_docker_definitions.GetSrcFolderList(dockerPaths)

	if len(srcFolderList) > 0 {
		for _, folder := range srcFolderList {
			// Creating missing project Docker folders
			if err := xfs.CreateFolder("ssh://" + folder.Path); err != nil {
				return err
			}
			// Adjust Docker folder permissions
			if folder.User != "" {
				resolvedUserId, err := resolveUserNameWithCache(folder.User)
				if err != nil {
					return fmt.Errorf("Unable to resolve user \"" + folder.User + "\" into a UID")
				}
				// Change user of folder
				if _, _, err := system.RemoteRun("chown "+strconv.Itoa(resolvedUserId)+":stackhead "+folder.Path, system.RemoteRunOpts{Sudo: true}); err != nil {
					return err
				}
			}
		}
	}

	// remove old hook files
	if _, _, err := system.RemoteRun("rm -rf "+dockerPaths.GetHooksDir(), system.RemoteRunOpts{}); err != nil {
		return err
	}
	if err := xfs.CreateFolder("ssh://" + dockerPaths.GetHooksDir()); err != nil {
		return err
	}

	// Collect new hooks
	var collectedHooks []container_docker_definitions.Hook
	for _, service := range system.Context.Project.Container.Services {
		if service.Hooks.ExecuteAfterSetup != "" {
			collectedHooks = append(collectedHooks, container_docker_definitions.Hook{
				Prefix:  "afterSetup",
				Src:     path.Join(system.Context.Project.ProjectDefinitionFolder, service.Hooks.ExecuteAfterSetup),
				Service: service.Name,
			})
		}
		if service.Hooks.ExecuteBeforeDestroy != "" {
			collectedHooks = append(collectedHooks, container_docker_definitions.Hook{
				Prefix:  "beforeDestroy",
				Src:     path.Join(system.Context.Project.ProjectDefinitionFolder, service.Hooks.ExecuteBeforeDestroy),
				Service: service.Name,
			})
		}
	}

	// Copy hook files
	for _, hook := range collectedHooks {
		hasHook, err := xfs.HasFile(hook.Src)
		if err != nil {
			return fmt.Errorf("Unable to validate hook file's existence: " + err.Error())
		}
		if !hasHook {
			return fmt.Errorf("Missing hook file \"" + hook.Src + "\"")
		}
		remoteHookFilePath := path.Join(dockerPaths.GetHooksDir(), hook.Service, hook.Prefix+"_"+path.Base(hook.Src))
		if err := xfs.CreateFolder("ssh://" + path.Dir(remoteHookFilePath)); err != nil {
			return err
		}
		if err := xfs.CopyFile(
			hook.Src,
			"ssh://"+remoteHookFilePath,
		); err != nil {
			return err
		}
	}

	composeYaml, err := docker_compose.BuildDockerCompose(system.Context.Project)
	if err != nil {
		return err
	}
	composeMap, err := composeYaml.Map()
	if err != nil {
		return err
	}

	dockerComposeResource := system.Resource{
		Type:      system.TypeFile,
		Operation: system.OperationCreate,
		Name:      "docker-compose.yaml",
	}

	oldComposeFilePath := ""
	var remoteComposeObjMap map[string]interface{}
	if system.Context.LatestDeployment != nil {
		dockerComposeFilePathOld, err := system.Context.LatestDeployment.GetResourcePath(dockerComposeResource)
		if err != nil {
			return err
		}
		oldComposeFilePath = dockerComposeFilePathOld
		hasRemoteFile, err := xfs.HasFile("ssh://" + oldComposeFilePath)
		if err != nil && hasRemoteFile {
			remoteComposeContent, err := xfs.ReadFile("ssh://" + oldComposeFilePath)
			if err != nil {
				return fmt.Errorf("unable to read remote docker-compose.yaml file from previous deployment: " + err.Error())
			}
			remoteComposeObj := docker_compose.DockerCompose{}
			if err := yaml.Unmarshal([]byte(remoteComposeContent), &remoteComposeObj); err != nil {
				return fmt.Errorf("unable to read remote docker-compose.yaml file from previous deployment: " + err.Error())
			}
			remoteComposeObjMap, err = remoteComposeObj.Map()
			if err != nil {
				return fmt.Errorf("unable to process remote docker-compose.yaml file from previous deployment: " + err.Error())
			}
		} else if err != nil && err.Error() != "file does not exist" {
			return fmt.Errorf("Unable to check state of remote docker-compose.yaml from previous deployment: " + err.Error())
		}
	}

	result := diff_docker_compose.DiffYaml(remoteComposeObjMap, composeMap)
	updateRequired, err := prepareUpdate(result)
	if err != nil {
		return err
	}

	if len(result.Diffs) == 0 && !updateRequired {
		// todo: pass to command task output
		//logger.Infoln("No changes to Docker Compose file found and Docker images are up-to-date. No need to update.")
		return nil
	}
	evaluateDiff(result)

	composeFileContent, err := composeYaml.String()
	if err != nil {
		return err
	}
	dockerComposeResource.Content = composeFileContent

	dockerComposeFilePathNew, _ := system.Context.CurrentDeployment.GetResourcePath(dockerComposeResource)

	system.Context.CurrentDeployment.ResourceGroups = append(system.Context.CurrentDeployment.ResourceGroups, system.ResourceGroup{
		Name:      "container-docker-" + system.Context.Project.Name + "-composefile",
		Resources: []system.Resource{dockerComposeResource},
	})

	var containerResources []system.Resource
	for serviceName, service := range composeYaml.Services {
		containerResources = append(containerResources, system.Resource{
			Type:        system.TypeContainer,
			Operation:   system.OperationCreate,
			ServiceName: serviceName,
			Name:        service.ContainerName,
			ImageName:   service.Image,
			Ports:       service.Ports,
		})
	}

	system.Context.CurrentDeployment.ResourceGroups = append(system.Context.CurrentDeployment.ResourceGroups, system.ResourceGroup{
		Name:      "container-docker-" + system.Context.Project.Name + "-containers",
		Resources: containerResources,
		ApplyResourceFunc: func() error {
			if oldComposeFilePath != "" {
				// Stop old Docker Compose containers
				// todo: allow using either docker-compose or "docker compose" whichever is available (prefer "docker compose")
				if _, stderr, err := system.RemoteRun("docker compose", system.RemoteRunOpts{Args: []string{"down"}, WorkingDir: path.Dir(oldComposeFilePath)}); err != nil {
					if stderr.Len() > 0 {
						return fmt.Errorf("Unable to stop old Docker containers: " + stderr.String())
					}
					return fmt.Errorf("Unable to stop old Docker containers: " + err.Error())
				}
			}

			// Start Docker Compose
			// todo: allow using either docker-compose or "docker compose" whichever is available (prefer "docker compose")
			if _, stderr, err := system.RemoteRun("docker compose", system.RemoteRunOpts{Args: []string{"up", "-d"}, WorkingDir: path.Dir(dockerComposeFilePathNew)}); err != nil {
				if stderr.Len() > 0 {
					return fmt.Errorf("Unable to start Docker containers: " + stderr.String())
				}
				return fmt.Errorf("Unable to start Docker containers: " + err.Error())
			}

			// Execute hooks
			if err := docker_system.ExecuteHook("afterSetup"); err != nil {
				return fmt.Errorf("After setup hook %s failed: ", err.Error())
			}
			return nil
		},
		RollbackResourceFunc: func() error {
			// Start old containers again
			if oldComposeFilePath != "" {
				// todo: allow using either docker-compose or "docker compose" whichever is available (prefer "docker compose")
				if _, stderr, err := system.RemoteRun("docker compose", system.RemoteRunOpts{Args: []string{"up", "-d"}, WorkingDir: path.Dir(oldComposeFilePath)}); err != nil {
					if stderr.Len() > 0 {
						return fmt.Errorf("Unable to stop Docker containers: " + stderr.String())
					}
					return fmt.Errorf("Unable to start old Docker containers: " + err.Error())
				}
			}

			// Stop Docker Compose
			// todo: allow using either docker-compose or "docker compose" whichever is available (prefer "docker compose")
			if _, stderr, err := system.RemoteRun("docker compose", system.RemoteRunOpts{Args: []string{"down"}, WorkingDir: path.Dir(dockerComposeFilePathNew)}); err != nil {
				if stderr.Len() > 0 {
					return fmt.Errorf("Unable to stop Docker containers: " + stderr.String())
				}
				return fmt.Errorf("Unable to stop new Docker containers: " + err.Error())
			}
			return nil
		},
	})

	return nil
}

func prepareUpdate(result diff_docker_compose.YamlDiffResult) (bool, error) {
	for _, registry := range system.Context.Project.Container.Registries {
		_, err := system.SimpleRemoteRun("docker", system.RemoteRunOpts{Args: []string{"login", "-u " + registry.Username, "-p " + registry.Password, registry.Url}, Confidential: true})
		if err != nil {
			return false, err
		}
	}

	updatedImages := false
	changedServices := result.GetStructure([]string{"services"})
	for _, structure := range changedServices.GetChildren() {
		// Look for new Docker image
		if structure.GetDiff().GetType() == diff_docker_compose.Unchanged || structure.GetDiff().GetType() == diff_docker_compose.Added {
			// Convert ValueNew into services object
			jsonStr, err := json.Marshal(structure.GetDiff().ValueNew)
			if err != nil {
				return false, err
			}
			var service docker_compose.Services
			if err := json.Unmarshal(jsonStr, &service); err != nil {
				return false, err
			}

			stdout, stderr, err := system.RemoteRun("docker", system.RemoteRunOpts{Args: []string{"pull", service.Image}})
			if err != nil {
				return false, fmt.Errorf("Unable to pull image from registry: " + stderr.String())
			}
			output := stdout.String()
			logger.Debugln(output)
			if strings.Contains(output, "Downloaded newer image for "+service.Image) {
				updatedImages = true
				// Image was downloaded
				digestMatch := regexp.MustCompile(`(?m)^Digest: (.*)$`).FindStringSubmatch(output)
				if len(digestMatch) > 1 {
					logger.Infoln("Downloaded newer image for " + service.Image + " (Digest " + digestMatch[1] + ")")
				} else {
					logger.Infoln("Downloaded newer image for " + service.Image)
				}
				// todo: log change to system
			}
		}
	}
	return updatedImages, nil
}

func evaluateDiff(result diff_docker_compose.YamlDiffResult) {
	if !result.HasChanged([]string{"services"}) {
		return
	}
	serviceStructure := result.GetStructure([]string{"services"})
	var addedServices []string
	var removedServices []string
	var modifiedServices []string

	for serviceName, service := range serviceStructure.GetChildren() {
		switch service.GetDiff().GetType() {
		case diff_docker_compose.Added:
			addedServices = append(addedServices, serviceName)
			break
		case diff_docker_compose.Removed:
			removedServices = append(removedServices, serviceName)
			break
		case diff_docker_compose.Changed:
			modifiedServices = append(modifiedServices, serviceName)
			break
		}
	}

	if len(removedServices) > 0 {
		logger.Infoln("Services locally removed/disabled:")
		for _, service := range removedServices {
			logger.Infoln("* " + service)
		}
		logger.Infoln("")
	}

	if len(addedServices) > 0 {
		logger.Infoln("Services locally added/enabled:")
		for _, service := range addedServices {
			logger.Infoln("* " + service)
		}
		logger.Infoln("")
	}

	if len(modifiedServices) > 0 {
		logger.Infoln("Services locally modified:")
		for _, service := range modifiedServices {
			logger.Infoln("* " + service)
		}
		logger.Infoln("")
	}
}
