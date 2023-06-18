package container_docker_definitions

import (
	"github.com/getstackhead/stackhead/project"
	"github.com/getstackhead/stackhead/system"
	"path"
)

type DockerPaths struct {
	DataDir       string
	DeploymentDir string
}

func GetDockerPaths() DockerPaths {
	return DockerPaths{
		DataDir:       path.Join(system.Context.Project.GetRuntimeDataDirectoryPath(), "container"),
		DeploymentDir: path.Join(system.Context.CurrentDeployment.GetPath(), "container"),
	}
}

func (p DockerPaths) GetHooksDir() string {
	return path.Join(p.DeploymentDir, "hooks")
}

func (p DockerPaths) getDataDir() string {
	return path.Join(p.DataDir, "data")
}

func (p DockerPaths) GetServiceDataDir(service project.ContainerService, volume project.ContainerServiceVolume) string {
	return p.getDataDir() + "/services/" + service.Name + "/" + volume.Src + "/"
}

func (p DockerPaths) GetGlobalDataDir(volume project.ContainerServiceVolume) string {
	return p.getDataDir() + "/global/" + volume.Src + "/"
}

type DockerVolumeInformation struct {
	Path string
	User string
}

func GetSrcFolderList(paths DockerPaths) []DockerVolumeInformation {
	var folders []DockerVolumeInformation

	// Container hooks location
	folders = append(folders, DockerVolumeInformation{
		Path: paths.GetHooksDir(),
		User: "",
	})

	for _, service := range system.Context.Project.Container.Services {
		for _, volume := range service.Volumes {
			// Collect local volumes
			if volume.Type == "local" {
				folders = append(folders, DockerVolumeInformation{
					Path: paths.GetServiceDataDir(service, volume),
					User: volume.User,
				})
			} else if volume.Type == "global" {
				folders = append(folders, DockerVolumeInformation{
					Path: paths.GetGlobalDataDir(volume),
					User: volume.User,
				})
			} else if volume.Type == "custom" {
				folders = append(folders, DockerVolumeInformation{
					Path: volume.Src,
					User: volume.User,
				})
			}
		}
	}

	return folders
}
