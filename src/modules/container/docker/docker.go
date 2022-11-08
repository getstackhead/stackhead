package container_docker

import (
	"embed"

	"github.com/getstackhead/stackhead/system"
)

type Module struct {
}

//go:embed templates
var templates embed.FS

func (Module) GetTemplates() *embed.FS {
	return &templates
}

func (Module) GetConfig() system.ModuleConfig {
	return system.ModuleConfig{
		Name: "docker",
		Type: "container",
		Terraform: system.ModuleTerraformConfig{
			Provider: system.ModuleTerraformConfigProvider{
				Vendor:             "kreuzwerker",
				Name:               "docker",
				Version:            "2.20.0",
				ResourceName:       "docker_container",
				Init:               "provider_init.tf.tmpl", // relative to "./templates/",
				ProviderPerProject: true,
			},
		},
	}
}
