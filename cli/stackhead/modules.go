package stackhead

import (
	"fmt"
	"strings"
)

const (
	ModuleWebserver = "stackhead_webserver"
	ModuleContainer = "stackhead_container"
)

func SplitModuleName(moduleName string) (string, string, string) {
	return ExtractVendor(moduleName), GetModuleType(moduleName), GetModuleBaseName(moduleName)
}

func GetModuleBaseName(moduleName string) string {
	moduleName = RemoveVendor(moduleName)
	moduleType := GetModuleType(moduleName)
	return strings.TrimPrefix(moduleName, moduleType)
}

func ExtractVendor(moduleName string) string {
	if !strings.ContainsRune(moduleName, '.') {
		return ""
	}
	var split = strings.Split(moduleName, ".")
	return split[0]
}

func RemoveVendor(moduleName string) string {
	if !strings.ContainsRune(moduleName, '.') {
		return moduleName
	}
	var split = strings.Split(moduleName, ".")
	return split[1]
}

func IsWebserverModule(moduleName string) bool {
	moduleName = RemoveVendor(moduleName)
	return strings.HasPrefix(moduleName, ModuleWebserver)
}

func IsContainerModule(moduleName string) bool {
	moduleName = RemoveVendor(moduleName)
	return strings.HasPrefix(moduleName, ModuleContainer)
}

func GetModuleType(moduleName string) string {
	if IsContainerModule(moduleName) {
		return ModuleContainer
	}
	if IsWebserverModule(moduleName) {
		return ModuleWebserver
	}
	return ""
}

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

	return vendorName + "." + moduleType + "_" + baseName, nil
}
