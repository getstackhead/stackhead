package plugin_portainer

import (
	"embed"
	container_docker "github.com/getstackhead/stackhead/modules/container/docker"

	"github.com/getstackhead/stackhead/system"
)

type Module struct {
}

// go:embed templates
var templates embed.FS

func (Module) Install(moduleSettings interface{}) error {
	// not implemented for modules of type "plugin"
	return nil
}

func (Module) GetTemplates() *embed.FS {
	return &templates
}

func (Module) GetConfig() system.ModuleConfig {
	config := container_docker.Module{}.GetConfig()
	config.Name = "portainer"
	return config
}
