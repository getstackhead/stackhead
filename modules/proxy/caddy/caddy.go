package proxy_caddy

import "github.com/getstackhead/stackhead/system"

type Module struct {
}

func (Module) GetConfig() system.ModuleConfig {
	return system.ModuleConfig{
		Name: "caddy",
		Type: "proxy",
		Terraform: system.ModuleTerraformConfig{
			Provider: system.ModuleTerraformConfigProvider{
				Vendor:       "getstackhead",
				Name:         "caddy",
				Version:      "1.0.1",
				ResourceName: "caddy_server_block",
			},
		},
	}
}
