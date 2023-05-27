package docker_system

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/phayes/freeport"

	"github.com/getstackhead/stackhead/project"
	"github.com/getstackhead/stackhead/system"
)

func GetPortMap(project *project.Project) (map[string]int, error) {
	dockerPortMap := map[string]int{}

	// find ports for running containers
	for _, service := range project.Container.Services {
		res, _, err := system.RemoteRun("docker", system.RemoteRunOpts{Args: []string{"port", ContainerName(project.Name, service.Name, system.Context.CurrentDeployment)}})
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
