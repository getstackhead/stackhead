package proxy_nginx

import (
	"text/template"

	"github.com/getstackhead/stackhead/system"
)

type NginxProxyModule struct {
}

func (NginxProxyModule) GetConfig() system.ModuleConfig {
	SnakeoilFullchainPath, SnakeoilPrivkeyPath := GetSnakeoilPaths()
	return system.ModuleConfig{
		Type: "proxy",
		Terraform: system.ModuleTerraformConfig{
			Provider: system.ModuleTerraformConfigProvider{
				Vendor:       "getstackhead",
				Name:         "nginx",
				Version:      "1.3.2",
				ResourceName: "nginx_server_block",
				Init:         "proxy/nginx/providers.tf.tmpl",
				InitFuncMap: template.FuncMap{
					"SnakeoilFullchainPath": func() string {
						return SnakeoilFullchainPath
					},
					"SnakeoilPrivkeyPath": func() string {
						return SnakeoilPrivkeyPath
					},
				},
			},
		},
	}
}
