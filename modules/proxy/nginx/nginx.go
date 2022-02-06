package proxy_nginx

import "github.com/getstackhead/stackhead/system"

type NginxProxyModule struct {
}

func (NginxProxyModule) GetConfig() system.ModuleConfig {
	return system.ModuleConfig{
		Type: "proxy",
		Terraform: system.ModuleTerraformConfig{
			Provider: system.ModuleTerraformConfigProvider{
				Vendor:       "getstackhead",
				Name:         "nginx",
				Version:      "1.3.2",
				ResourceName: "nginx_server_block",
			},
		},
	}
}
