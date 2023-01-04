package proxy_nginx

import (
	"embed"

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
	return system.ModuleConfig{
		Name: "nginx",
		Type: "proxy",
	}
}
