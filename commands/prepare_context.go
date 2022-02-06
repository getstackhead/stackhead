package commands

import (
	container_docker "github.com/getstackhead/stackhead/modules/container/docker"
	proxy_nginx "github.com/getstackhead/stackhead/modules/proxy/nginx"
	"github.com/getstackhead/stackhead/project"
	"github.com/getstackhead/stackhead/system"
)

func PrepareContext(host string, action string, projectDefinition *project.Project) {
	system.InitializeContext(host, action, projectDefinition)
	system.ContextSetProxyModule(proxy_nginx.NginxProxyModule{})
	system.ContextSetContainerModule(container_docker.DockerContainerModule{})
}
