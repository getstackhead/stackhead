package docker_system

import (
	"fmt"
	"path"
	"strings"

	xfs "github.com/saitho/golang-extended-fs/v2"

	container_docker_definitions "github.com/getstackhead/stackhead/modules/container/docker/definitions"
	"github.com/getstackhead/stackhead/system"
)

func ExecuteHook(hookName string) error {
	hooksDir := container_docker_definitions.GetDockerPaths().GetHooksDir()
	info, err := xfs.ListFolders("ssh://" + hooksDir)
	if err != nil {
		return err
	}
	var files []struct {
		File    string
		Service string
	}
	for _, fileInfo := range info { // first level: service directories
		if !fileInfo.IsDir() {
			continue
		}
		hookFiles, err := xfs.ListFolders("ssh://" + path.Join(hooksDir, fileInfo.Name()))
		if err != nil {
			return err
		}
		for _, hookInfo := range hookFiles { // first level: service directories
			if hookInfo.IsDir() || !strings.HasPrefix(hookInfo.Name(), hookName+"_") {
				continue
			}
			files = append(files, struct {
				File    string
				Service string
			}{File: hookInfo.Name(), Service: fileInfo.Name()})
		}
	}

	for _, file := range files {
		filePath := path.Join(hooksDir, file.Service, file.File)

		// copy file onto container and run it.....
		containerLocation := path.Join("/", file.File)
		containerName := ContainerName(system.Context.Project.Name, file.Service, system.Context.CurrentDeployment)
		_, err := system.SimpleRemoteRun("docker", system.RemoteRunOpts{
			Args: []string{
				"cp",
				filePath,
				containerName + ":" + containerLocation,
			},
			WorkingDir: system.Context.CurrentDeployment.GetPath(),
		})
		if err != nil {
			return fmt.Errorf("Unable to copy file %s to container %s: \"%s\"", file.File, containerName, err.Error())
		}
		_, err = system.SimpleRemoteRun("docker", system.RemoteRunOpts{
			Args: []string{
				"exec",
				"-u 0",
				containerName,
				"chmod +x " + containerLocation,
			},
			WorkingDir: system.Context.CurrentDeployment.GetPath(),
		})
		if err != nil {
			return fmt.Errorf("Unable to copy file %s to container %s: \"%s\"", file.File, containerName, err.Error())
		}

		_, err = system.SimpleRemoteRun("docker", system.RemoteRunOpts{
			Args: []string{
				"exec",
				containerName,
				containerLocation,
			},
			WorkingDir: system.Context.CurrentDeployment.GetPath(),
		})
		if err != nil {
			return fmt.Errorf("Unable to run %s on container %s: \"%s\"", containerLocation, containerName, err.Error())
		}
	}

	return nil
}
