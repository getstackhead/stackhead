package docker_compose

import (
	"fmt"
	"path"
	"strings"

	"golang.org/x/exp/slices"

	container_docker_definitions "github.com/getstackhead/stackhead/modules/container/docker/definitions"
	docker_system "github.com/getstackhead/stackhead/modules/container/docker/system"
	"github.com/getstackhead/stackhead/project"
	"github.com/getstackhead/stackhead/system"
)

func BuildDockerCompose(project *project.Project) (DockerCompose, error) {
	dockerPaths := container_docker_definitions.GetDockerPaths()
	compose := DockerCompose{
		Version:  "2.4",
		Networks: map[string]Network{"stackhead-network-" + project.Name: {}},
		Services: map[string]Services{},
		Volumes:  map[string]Volume{},
	}
	dockerPortMap, err := docker_system.GetPortMap(system.Context.Project)
	if err != nil {
		return compose, err
	}

	for _, service := range project.Container.Services {
		addService(&compose, project, service, dockerPortMap)
	}

	for _, service := range project.Container.Services {
		for _, volume := range service.Volumes {
			if !slices.Contains([]string{"global", "custom", "local"}, volume.Type) {
				continue
			}
			vol := Volume{}
			serviceName := service.Name
			vol.DriverOpts.Type = "none"
			vol.DriverOpts.O = "bind"
			if volume.Type == "local" {
				vol.DriverOpts.Device = path.Join(dockerPaths.GetServiceDataDir(service, volume))
			} else if volume.Type == "global" {
				vol.DriverOpts.Device = path.Join(dockerPaths.GetGlobalDataDir(volume))
			} else if volume.Type == "custom" {
				vol.DriverOpts.Device = volume.Src
			}
			compose.Volumes[GetVolumeSrcKey(project.Name, serviceName, volume)] = vol
		}
	}

	return compose, nil
}

func addService(compose *DockerCompose, project *project.Project, service project.ContainerService, dockerPortMap map[string]int) {
	var volumes []ServiceVolume
	var dependOnServices []string
	for _, volume := range service.Volumes {
		serviceName := service.Name
		volumes = append(volumes, ServiceVolume{
			Type:     "volume",
			Source:   GetVolumeSrcKey(project.Name, serviceName, volume),
			Target:   volume.Dest,
			ReadOnly: volume.Mode == "ro",
		})
	}

	var volumesFrom []string
	for _, volumeFrom := range service.VolumesFrom {
		split := strings.SplitN(volumeFrom, ":", 2)
		containerName := split[0]
		dependOnServices = append(dependOnServices, containerName)
		volumesFrom = append(volumesFrom, volumeFrom)
	}

	var ports []string
	for _, domain := range project.Domains {
		for _, expose := range domain.Expose {
			if expose.Service != service.Name {
				continue
			}
			portMapKey := fmt.Sprintf("%s-%d", expose.Service, expose.InternalPort)
			portString := fmt.Sprintf("%d:%d", dockerPortMap[portMapKey], expose.InternalPort)
			if !slices.Contains(ports, portString) {
				ports = append(ports, portString)
			}
		}
	}

	compose.Services[service.Name] = Services{
		ContainerName: docker_system.ContainerName(project.Name, service.Name),
		Image:         service.Image,
		Restart:       "unless-stopped",
		Labels:        map[string]string{"stackhead.project": project.Name},
		User:          service.User,
		Networks:      map[string]ServiceNetwork{"stackhead-network-" + project.Name: {Aliases: []string{service.Name}}},
		Volumes:       volumes,
		Ports:         ports,
		Environment:   service.Environment,
		DependsOn:     dependOnServices,
		VolumesFrom:   volumesFrom,
	}
}
