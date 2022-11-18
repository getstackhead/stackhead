package commands

import (
	"github.com/knadh/koanf/maps"
	xfs "github.com/saitho/golang-extended-fs/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/getstackhead/stackhead/config"
	"github.com/getstackhead/stackhead/project"
	"github.com/getstackhead/stackhead/system"

	container_docker "github.com/getstackhead/stackhead/modules/container/docker"
	dns_cloudflare "github.com/getstackhead/stackhead/modules/dns/cloudflare"
	plugin_portainer "github.com/getstackhead/stackhead/modules/plugin/portainer"
	proxy_caddy "github.com/getstackhead/stackhead/modules/proxy/caddy"
	proxy_nginx "github.com/getstackhead/stackhead/modules/proxy/nginx"
)

func EnforceSimpleValueOption(name string, value map[string]interface{}) {
	stringMap := viper.GetStringMap(name)
	for mapKey, mapValue := range value {
		stringMap[mapKey] = mapValue
		logger.Warnf("Enforcing setting \"%s.%s=%v\" via remote server configuration", name, mapKey, mapValue)
	}
	viper.Set(name, stringMap)
}

func EnforceNestedValueOption(name string, value map[string]interface{}) {
	modulesConfigMap := viper.GetStringMap(name)
	for moduleName, moduleConfig := range cast.ToStringMap(value) {
		remoteSettings := cast.ToStringMap(moduleConfig)
		logMessage := "Enforcing module settings for \"" + moduleName + "\" via remote server configuration"
		// todo: uncomment when we can sanitize secrets in CLI output
		//for k, v := range remoteSettings {
		//	logMessage += fmt.Sprintf("\n"+k+"=%v", v)
		//}
		logger.Warn(logMessage)
		newMap := cast.ToStringMap(modulesConfigMap[moduleName])
		maps.Merge(remoteSettings, newMap)
		modulesConfigMap[moduleName] = newMap
	}
	viper.Set(name, modulesConfigMap)
}

func PrepareContext(host string, action string, projectDefinition *project.Project) {
	system.InitializeContext(host, action, projectDefinition)

	// Check remote config on server
	hasFile, _ := xfs.HasFile("ssh://" + config.GetServerConfigFilePath())
	if hasFile {
		fileContent, err := xfs.ReadFile("ssh://" + config.GetServerConfigFilePath())
		if err != nil {
			panic("Found remote config file but was unable to read it: " + err.Error())
		}
		var c map[string]interface{}
		if err = yaml.Unmarshal([]byte(fileContent), &c); err != nil {
			panic("Found remote config file but was unable to parse it: " + err.Error())
		}

		// Enforce remote configurations on viper
		EnforceSimpleValueOption("modules", cast.ToStringMap(c["modules"]))
		EnforceNestedValueOption("modules_config", cast.ToStringMap(c["modules_config"]))
	}

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
	dnsNames := viper.GetStringMapStringSlice("modules")["dns"]
	for _, dnsName := range dnsNames {
		switch dnsName {
		case "cloudflare":
			system.ContextAddDnsModule(dns_cloudflare.Module{})
		}
	}

	// set plugin
	pluginNames := viper.GetStringMapStringSlice("modules")["plugins"]
	for _, pluginName := range pluginNames {
		switch pluginName {
		case "portainer":
			system.ContextAddPluginModule(plugin_portainer.Module{})
		}
	}

	// todo: validate modules
}
