package docker_system

import (
	"fmt"
	"github.com/getstackhead/stackhead/system"
)

func ContainerName(projectName string, serviceName string, deployment system.Deployment) string {
	return fmt.Sprintf("stackhead-%s-%s-v%d", projectName, serviceName, deployment.Version)
}

func NetworkName(projectName string, deployment system.Deployment) string {
	return fmt.Sprintf("stackhead-network-%s-v%d", projectName, deployment.Version)
}
