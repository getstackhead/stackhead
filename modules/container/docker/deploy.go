package container_docker

import (
	"fmt"
	"path"
	"reflect"
	"regexp"
	"strconv"
	"text/template"

	"github.com/phayes/freeport"
	xfs "github.com/saitho/golang-extended-fs/v2"

	container_docker_definitions "github.com/getstackhead/stackhead/modules/container/docker/definitions"
	"github.com/getstackhead/stackhead/project"
	"github.com/getstackhead/stackhead/system"
)

func GetDockerPaths() container_docker_definitions.DockerPaths {
	return container_docker_definitions.DockerPaths{
		BaseDir: system.Context.Project.GetRuntimeDataDirectoryPath() + "/container",
	}
}

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

func getPortMap(project *project.Project) (map[string]int, error) {
	dockerPortMap := map[string]int{}

	// find ports for running containers
	for _, service := range project.Container.Services {
		res, _, err := system.RemoteRun("docker", "port", "stackhead-"+project.Name+"-"+service.Name)
		if err == nil { // ignore error (container not running)
			// e.g. 80/tcp -> 0.0.0.0:49155
			re := regexp.MustCompile(`(?P<Internal>\d+)\/tcp -> 0\.0\.0\.0:(?P<External>\d+)`)
			matches := re.FindAllStringSubmatch(res.String(), -1)
			for _, match := range matches {
				externalPort, _ := strconv.Atoi(match[re.SubexpIndex("External")])
				dockerPortMap[service.Name+"-"+match[re.SubexpIndex("Internal")]] = externalPort
			}
		}
	}

	// determine ports for missing containers
	missingPortServices := []string{}
	for _, domain := range project.Domains {
		for _, expose := range domain.Expose {
			mapKey := expose.Service + "-" + strconv.Itoa(expose.InternalPort)
			if _, ok := dockerPortMap[mapKey]; !ok {
				missingPortServices = append(missingPortServices, mapKey)
			}
		}
	}
	ports, err := freeport.GetFreePorts(len(missingPortServices))
	if err != nil {
		return nil, fmt.Errorf("unable to determine free ports: " + err.Error())
	}
	for i := range missingPortServices {
		dockerPortMap[missingPortServices[i]] = ports[i]
	}
	return dockerPortMap, nil
}

func (m Module) Deploy(modulesSettings interface{}) error {
	project := system.Context.Project

	// Build src folder list
	srcFolderList := container_docker_definitions.GetSrcFolderList(GetDockerPaths())

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
				if _, _, err := system.RemoteRun("sudo chown " + strconv.Itoa(resolvedUserId) + ":stackhead " + folder.Path); err != nil {
					return err
				}
			}
		}
	}

	// remove old hook files
	if _, _, err := system.RemoteRun("rm -rf " + GetDockerPaths().GetHooksDir()); err != nil {
		return err
	}
	if err := xfs.CreateFolder("ssh://" + GetDockerPaths().GetHooksDir()); err != nil {
		return err
	}

	// Collect new hooks
	var collectedHooks []container_docker_definitions.Hook
	for _, service := range system.Context.Project.Container.Services {
		if service.Hooks.ExecuteAfterSetup != "" {
			collectedHooks = append(collectedHooks, container_docker_definitions.Hook{
				Prefix: "afterSetup",
				Src:    path.Join(system.Context.Project.ProjectDefinitionFolder, service.Hooks.ExecuteAfterSetup),
			})
		}
		if service.Hooks.ExecuteBeforeDestroy != "" {
			collectedHooks = append(collectedHooks, container_docker_definitions.Hook{
				Prefix: "beforeDestroy",
				Src:    path.Join(system.Context.Project.ProjectDefinitionFolder, service.Hooks.ExecuteBeforeDestroy),
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
		remoteHookFilePath := "ssh://" + path.Join(GetDockerPaths().GetHooksDir(), hook.Prefix+"_"+path.Base(hook.Src))
		if err := xfs.CopyFile(
			hook.Src,
			remoteHookFilePath,
		); err != nil {
			return err
		}
		if err := xfs.Chmod(remoteHookFilePath, 0755); err != nil {
			return err
		}
	}

	// Generate Terraform Docker configuration file
	var funcMap = template.FuncMap{
		"sanitize_volume": func(s string) string {
			var re = regexp.MustCompile(`[^\w]`)
			return re.ReplaceAllString(s, "_")
		},
		// Container specific
		"TF_replace": func(input string, projectName string) string {
			// Replace Docker service name variables
			// Example: TF_replace "$DOCKER_SERVICE_NAME['0'] - $DOCKER_SERVICE_NAME['1']" "myproject"
			// Result: ${docker_container.stackhead-myproject-0.name} - ${docker_container.stackhead-myproject-1.name}

			var re = regexp.MustCompile("\\$DOCKER_SERVICE_NAME['(.*)']")
			resource := m.GetConfig().Terraform.Provider.ResourceName
			return re.ReplaceAllString(input, "${"+resource+".stackhead-"+projectName+"-\\1.name}")
		},
	}

	dockerPortMap, err := getPortMap(project)
	if err != nil {
		return fmt.Errorf("unable to determine free ports: " + err.Error())
	}

	data := map[string]interface{}{
		"Context":       system.Context,
		"DockerPaths":   GetDockerPaths(),
		"DockerPortMap": dockerPortMap,
	}
	dockerTf, err := system.RenderModuleTemplate(
		templates,
		"project.tf.tmpl",
		data,
		funcMap)
	if err != nil {
		return err
	}

	_ = xfs.WriteFile("ssh://"+project.GetTerraformDirectoryPath()+"/docker.tf", dockerTf)
	return xfs.WriteFile("ssh://"+project.GetTerraformDirectoryPath()+"/docker.tf", dockerTf)
}
