package commands

import (
	"github.com/spf13/viper"

	container_docker "github.com/getstackhead/stackhead/modules/container/docker"
	proxy_caddy "github.com/getstackhead/stackhead/modules/proxy/caddy"
	proxy_nginx "github.com/getstackhead/stackhead/modules/proxy/nginx"
	"github.com/getstackhead/stackhead/project"
	"github.com/getstackhead/stackhead/system"
)

func PrepareContext(host string, action string, projectDefinition *project.Project) {
	system.InitializeContext(host, action, projectDefinition)

	// set proxy
	switch viper.GetStringMapString("modules")["proxy"] {
	case "nginx":
		system.ContextSetProxyModule(proxy_nginx.Module{})
	default: // use Caddy as default
		system.ContextSetProxyModule(proxy_caddy.Module{})
	}

	// set container
	switch viper.GetStringMapString("modules")["container"] {
	default: // use Docker as default
		system.ContextSetContainerModule(container_docker.Module{})
	}

	// set DNS
	dnsNames := viper.GetStringMap("modules")["dns"].([]interface{})
	for _, dnsName := range dnsNames {
		switch dnsName {
		// todo: add available DNS modules here
		}
	}

	// set plugin
	pluginNames := viper.GetStringMapStringSlice("modules")["plugins"]
	for _, pluginName := range pluginNames {
		switch pluginName {
		// todo: add available plugin modules here
		}
	}
}
