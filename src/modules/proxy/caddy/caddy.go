package proxy_caddy

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
