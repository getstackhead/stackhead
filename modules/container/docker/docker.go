package container_docker

import (
	"embed"

	"github.com/getstackhead/stackhead/system"
)

type Module struct {
}

func (Module) GetTemplates() *embed.FS {
	return nil
}

func (Module) GetConfig() system.ModuleConfig {
	return system.ModuleConfig{
		Name: "docker",
		Type: "container",
	}
}
