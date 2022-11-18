package dns_cloudflare

import (
	"embed"

	"github.com/getstackhead/stackhead/system"
)

type Module struct {
}

// go:embed templates
var templates embed.FS

func (Module) Install(moduleSettings interface{}) error {
	// not implemented for modules of type "dns"
	return nil
}

func (Module) GetTemplates() *embed.FS {
	return &templates
}

func (Module) GetConfig() system.ModuleConfig {
	return system.ModuleConfig{
		Name: "cloudflare",
		Type: "dns",
	}
}
