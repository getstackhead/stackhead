package stackhead

import (
	"fmt"
	"strings"
)

const (
	// ModuleWebserver is the string that identifies a package name as webserver package
	ModuleWebserver = "stackhead_webserver"
	// ModuleContainer is the string that identifies a package name as container package
	ModuleContainer = "stackhead_container"
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

// GetModuleType returns the module type for the given module according its name
// Will return the values of ModuleContainer and ModuleWebserver constants.
func GetModuleType(moduleName string) string {
	if IsContainerModule(moduleName) {
		return ModuleContainer
	}
	if IsWebserverModule(moduleName) {
		return ModuleWebserver
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
