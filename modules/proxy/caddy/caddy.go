package proxy_caddy

import "github.com/getstackhead/stackhead/system"

type CaddyProxyModule struct {
}

func (CaddyProxyModule) GetConfig() system.ModuleConfig {
	return system.ModuleConfig{
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
