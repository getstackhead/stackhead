package stackhead

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const (
	// ModuleWebserver is the string that identifies a package name as webserver package
	ModuleWebserver = "stackhead_webserver"
	// ModuleContainer is the string that identifies a package name as container package
	ModuleContainer = "stackhead_container"
	// ModulePlugin is the string that identifies a package name as plugin package
	ModulePlugin = "stackhead_plugin"
	// ModuleDNS is the string that identifies a package name as dns package
	ModuleDNS = "stackhead_dns"
)

// SplitModuleName splits a given module name into vendor, module type and base name
func SplitModuleName(moduleName string) (string, string, string) {
	return ExtractVendor(moduleName), GetModuleType(moduleName), GetModuleBaseName(moduleName)
}

// GetModuleBaseName returns the base name of a given module name (e.g. nginx in getstackhead.stackhead_webserver_nginx)
func GetModuleBaseName(moduleName string) string {
	moduleName = RemoveVendor(moduleName)
	moduleType := GetModuleType(moduleName)
	return strings.TrimPrefix(moduleName, moduleType+"_")
}

// ExtractVendor returns the vendor name of a given module name (e.g. getstackhead in getstackhead.stackhead_webserver_nginx)
func ExtractVendor(moduleName string) string {
	if !strings.ContainsRune(moduleName, '.') {
		return ""
	}
	var split = strings.Split(moduleName, ".")
	return split[0]
}

// RemoveVendor removes the vendor name from a given module name (e.g. getstackhead.stackhead_webserver_nginx => stackhead_webserver_nginx)
func RemoveVendor(moduleName string) string {
	if !strings.ContainsRune(moduleName, '.') {
		return moduleName
	}
	var split = strings.Split(moduleName, ".")
	return split[1]
}

// IsWebserverModule checks if the given module is a webserver module based on its name
func IsWebserverModule(moduleName string) bool {
	moduleName = RemoveVendor(moduleName)
	return strings.HasPrefix(moduleName, ModuleWebserver)
}

// IsContainerModule checks if the given module is a container module based on its name
func IsContainerModule(moduleName string) bool {
	moduleName = RemoveVendor(moduleName)
	return strings.HasPrefix(moduleName, ModuleContainer)
}

// IsPluginModule checks if the given module is a plugin module based on its name
func IsPluginModule(moduleName string) bool {
	moduleName = RemoveVendor(moduleName)
	return strings.HasPrefix(moduleName, ModulePlugin)
}

// IsDNSModule checks if the given module is a dns module based on its name
func IsDNSModule(moduleName string) bool {
	moduleName = RemoveVendor(moduleName)
	return strings.HasPrefix(moduleName, ModuleDNS)
}

// GetModuleType returns the module type for the given module according its name
// Will return the values of ModulePlugin, ModuleContainer and ModuleWebserver constants.
func GetModuleType(moduleName string) string {
	if IsContainerModule(moduleName) {
		return ModuleContainer
	}
	if IsPluginModule(moduleName) {
		return ModulePlugin
	}
	if IsWebserverModule(moduleName) {
		return ModuleWebserver
	}
	if IsDNSModule(moduleName) {
		return ModuleDNS
	}
	return ""
}

// AutoCompleteModuleName will try to auto-complete module names
// if no vendor is given, it will use "getstackhead"
// if no module type is given, it will use the given target type
func AutoCompleteModuleName(moduleNameFragment string, targetType string) (string, error) {
	var vendorName, moduleType, baseName = SplitModuleName(moduleNameFragment)
	if moduleType != "" && moduleType != targetType {
		return "", fmt.Errorf("invalid module name")
	}

	if vendorName == "" {
		vendorName = "getstackhead"
	}
	if moduleType == "" {
		moduleType = targetType
	}

	return vendorName + "." + strings.Join([]string{moduleType, baseName}, "_"), nil
}

func GetWebserverModule() (string, error) {
	var module = viper.GetString("modules.webserver")
	if len(module) == 0 {
		module = "getstackhead.stackhead_webserver_nginx"
	}
	return AutoCompleteModuleName(module, ModuleWebserver)
}

func GetContainerModule() (string, error) {
	var module = viper.GetString("modules.container")
	if len(module) == 0 {
		module = "getstackhead.stackhead_container_docker"
	}
	return AutoCompleteModuleName(module, ModuleContainer)
}

func GetDNSModules() ([]string, error) {
	var plugins = viper.GetStringSlice("modules.dns")
	var modules []string
	if len(plugins) > 0 {
		for _, plugin := range plugins {
			moduleName, err := AutoCompleteModuleName(plugin, ModuleDNS)
			if err != nil {
				return []string{}, err
			}
			modules = append(modules, moduleName)
		}
	}
	return modules, nil
}

func GetPluginModules() ([]string, error) {
	var plugins = viper.GetStringSlice("modules.plugins")
	var modules []string
	if len(plugins) > 0 {
		for _, plugin := range plugins {
			moduleName, err := AutoCompleteModuleName(plugin, ModulePlugin)
			if err != nil {
				return []string{}, err
			}
			modules = append(modules, moduleName)
		}
	}
	return modules, nil
}
