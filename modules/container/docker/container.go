package container_docker

import "github.com/getstackhead/stackhead/system"

type Module struct {
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
				Init:               "container/docker/provider_init.tf.tmpl", // relative to "templates/modules/",
				ProviderPerProject: true,
			},
		},
	}
}
