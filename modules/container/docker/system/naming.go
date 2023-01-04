package docker_system

import "fmt"

func ContainerName(projectName string, serviceName string) string {
	return fmt.Sprintf("stackhead-%s-%s", projectName, serviceName)
}
