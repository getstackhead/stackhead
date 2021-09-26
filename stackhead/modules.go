package stackhead

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const (
	// PluginProxy is the string that identifies a package name as webserver package
	PluginProxy = "stackhead_webserver"
	// PluginContainer is the string that identifies a package name as container package
	PluginContainer = "stackhead_container"
	// PluginApplication is the string that identifies a package name as application package
	PluginApplication = "stackhead_application"
	// PluginDNS is the string that identifies a package name as dns package
	PluginDNS = "stackhead_dns"
)

// SplitPluginName splits a given module name into vendor, module type and base name
func SplitPluginName(pluginName string) (string, string, string) {
	return ExtractVendor(pluginName), GetModuleType(pluginName), GetModuleBaseName(pluginName)
}

// GetModuleBaseName returns the base name of a given module name (e.g. nginx in getstackhead.stackhead_webserver_nginx)
func GetModuleBaseName(pluginName string) string {
	pluginName = RemoveVendor(pluginName)
	pluginType := GetModuleType(pluginName)
	return strings.TrimPrefix(pluginName, pluginType+"_")
}

// ExtractVendor returns the vendor name of a given module name (e.g. getstackhead in getstackhead.stackhead_webserver_nginx)
func ExtractVendor(pluginName string) string {
	if !strings.ContainsRune(pluginName, '.') {
		return ""
	}
	var split = strings.Split(pluginName, ".")
	return split[0]
}

// RemoveVendor removes the vendor name from a given module name (e.g. getstackhead.stackhead_webserver_nginx => stackhead_webserver_nginx)
func RemoveVendor(pluginName string) string {
	if !strings.ContainsRune(pluginName, '.') {
		return pluginName
	}
	var split = strings.Split(pluginName, ".")
	return split[1]
}

// IsProxyPlugin checks if the given module is a webserver module based on its name
func IsProxyPlugin(pluginName string) bool {
	pluginName = RemoveVendor(pluginName)
	return strings.HasPrefix(pluginName, PluginProxy)
}

// IsContainerPlugin checks if the given module is a container module based on its name
func IsContainerPlugin(pluginName string) bool {
	pluginName = RemoveVendor(pluginName)
	return strings.HasPrefix(pluginName, PluginContainer)
}

// IsApplicationPlugin checks if the given module is an application module based on its name
func IsApplicationPlugin(pluginName string) bool {
	pluginName = RemoveVendor(pluginName)
	return strings.HasPrefix(pluginName, PluginApplication)
}

// IsDNSPlugin checks if the given module is a dns module based on its name
func IsDNSPlugin(pluginName string) bool {
	pluginName = RemoveVendor(pluginName)
	return strings.HasPrefix(pluginName, PluginDNS)
}

// GetModuleType returns the module type for the given module according its name
// Will return the values of PluginMisc, PluginContainer and PluginProxy constants.
func GetModuleType(pluginName string) string {
	if IsContainerPlugin(pluginName) {
		return PluginContainer
	}
	if IsApplicationPlugin(pluginName) {
		return PluginApplication
	}
	if IsProxyPlugin(pluginName) {
		return PluginProxy
	}
	if IsDNSPlugin(pluginName) {
		return PluginDNS
	}
	return ""
}

// AutoCompletePluginName will try to auto-complete module names
// if no vendor is given, it will use "getstackhead"
// if no module type is given, it will use the given target type
func AutoCompletePluginName(pluginNameFragment string, targetType string) (string, error) {
	var vendorName, pluginType, baseName = SplitPluginName(pluginNameFragment)
	if pluginType != "" && pluginType != targetType {
		return "", fmt.Errorf("invalid module name")
	}

	if vendorName == "" {
		vendorName = "getstackhead"
	}
	if pluginType == "" {
		pluginType = targetType
	}

	return vendorName + "." + strings.Join([]string{pluginType, baseName}, "_"), nil
}

func GetProxyPlugin() string {
	var module = viper.GetString("modules.proxy")
	if len(module) == 0 {
		module = "github.com/getstackhead/plugin-proxy-nginx"
	}
	return module
}

func GetContainerPlugin() string {
	var module = viper.GetString("modules.container")
	if len(module) == 0 {
		module = "github.com/getstackhead/plugin-container-docker"
	}
	return module
}

func GetDNSPlugins() []string {
	return viper.GetStringSlice("modules.dns")
}

func GetApplicationPlugins() []string {
	return viper.GetStringSlice("modules.applications")
}
