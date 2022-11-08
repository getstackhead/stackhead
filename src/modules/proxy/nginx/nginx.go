package proxy_nginx

import (
	"embed"
	"text/template"

	"github.com/getstackhead/stackhead/system"
)

type Module struct {
}

//go:embed templates templates/**/*
var templates embed.FS

func (Module) GetTemplates() *embed.FS {
	return &templates
}

func (Module) GetConfig() system.ModuleConfig {
	SnakeoilFullchainPath, SnakeoilPrivkeyPath := GetSnakeoilPaths()
	return system.ModuleConfig{
		Name: "nginx",
		Type: "proxy",
		Terraform: system.ModuleTerraformConfig{
			Provider: system.ModuleTerraformConfigProvider{
				Vendor:       "getstackhead",
				Name:         "nginx",
				Version:      "1.3.2",
				ResourceName: "nginx_server_block",
				Init:         "providers.tf.tmpl",
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
