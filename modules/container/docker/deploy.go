package container_docker

import (
	"fmt"
	"path"
	"reflect"
	"regexp"
	"strconv"
	"text/template"

	xfs "github.com/saitho/golang-extended-fs/v2"

	container_docker_definitions "github.com/getstackhead/stackhead/modules/container/docker/definitions"
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
				var resolvedUserId int
				var resolvedGroupId int
				var err error
				if reflect.TypeOf(folder.User).String() == "string" { // if userId is a string, resolve it to uid
					resolvedUserId, err = system.ResolveRemoteUserIntoUid(folder.User)
				} else {
					resolvedUserId, err = strconv.Atoi(folder.User)
				}
				if err != nil {
					return fmt.Errorf("Unable to resolve user \"" + folder.User + "\" into a UID")
				}

				resolvedGroupId, err = system.ResolveRemoteGroupIntoGid("stackhead")
				if err != nil {
					return fmt.Errorf("Unable to resolve group \"stackhead\" into a GID. The StackHead setup seems to be incomplete.")
				}
				if err := xfs.Chown("ssh://"+folder.Path, resolvedUserId, resolvedGroupId); err != nil {
					return err
				}
			}
		}
	}

	// remove old hook files
	if _, _, err := system.RemoteRun("rm -rf " + GetDockerPaths().GetHooksDir()); err != nil {
		return err
	}

	// Collect new hooks
	var collectedHooks []container_docker_definitions.Hook
	for _, service := range system.Context.Project.Container.Services {
		if service.Hooks.ExecuteAfterSetup != "" {
			collectedHooks = append(collectedHooks, container_docker_definitions.Hook{
				Prefix: "afterSetup",
				Src:    service.Hooks.ExecuteAfterSetup,
			})
		}
		if service.Hooks.ExecuteBeforeDestroy != "" {
			collectedHooks = append(collectedHooks, container_docker_definitions.Hook{
				Prefix: "beforeDestroy",
				Src:    service.Hooks.ExecuteBeforeDestroy,
			})
		}
	}

	// Copy hook files
	for _, hook := range collectedHooks {
		hasHook, err := xfs.HasFile(hook.Src)
		fmt.Println(err)
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
	fmt.Println("Copy hook files")
	fmt.Println(collectedHooks)

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

	data := map[string]interface{}{
		"Context":     system.Context,
		"DockerPaths": container_docker_definitions.DockerPaths{},
	}
	dockerTf, err := system.RenderModuleTemplate(
		templates,
		"project.tf.tmpl",
		data,
		funcMap)
	if err != nil {
		return err
	}
	return xfs.WriteFile("ssh://"+project.GetTerraformDirectoryPath()+"/docker.tf", dockerTf)
}
