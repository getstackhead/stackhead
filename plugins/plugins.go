package plugins

import (
	"fmt"
	"github.com/getstackhead/stackhead/system"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/robertkrimen/otto"
	"gopkg.in/yaml.v3"

	"github.com/getstackhead/stackhead/config"
	"github.com/getstackhead/stackhead/pluginlib"
	"github.com/getstackhead/stackhead/stackhead"
)

type PluginProgram struct {
	OriginalPath string
	Source       string
}

func (p PluginProgram) Run() error {
	vm := otto.New()
	vm.Set("Stackhead", StackheadPluginFuncs)
	vm.Set("Project", system.Context.Project)
	vm.Set("__dirname", path.Dir(p.OriginalPath))
	vm.Set("__filename", p.OriginalPath)
	_, err := vm.Run(p.Source)
	return err
}

type Plugin struct {
	Name           string
	Path           string
	Config         *pluginlib.PluginConfig
	InitProgram    *PluginProgram
	SetupProgram   *PluginProgram
	DeployProgram  *PluginProgram
	DestroyProgram *PluginProgram
}

func SplitPluginPath(modulePath string) (string, string) {
	moduleName := modulePath
	moduleVersion := "main"
	lastInd := strings.LastIndex(modulePath, "@")
	if lastInd != -1 {
		moduleName = modulePath[:lastInd]
		moduleVersion = modulePath[lastInd+1:]
	}
	return moduleName, moduleVersion
}

func CollectPluginPaths() []string {
	var modules []string

	modules = append(modules, stackhead.GetProxyPlugin())
	//modules = append(modules, stackhead.GetContainerPlugin())
	//modules = append(modules, stackhead.GetDNSPlugins()...)
	//modules = append(modules, stackhead.GetApplicationPlugins()...)

	return modules
}

func LoadPlugins() ([]*Plugin, error) {
	var plugins []*Plugin
	pluginDir, err := config.GetPluginDir()
	if err != nil {
		return nil, err
	}
	for _, modulePath := range CollectPluginPaths() {
		moduleName, _ := SplitPluginPath(modulePath)
		pluginInstance, err := LoadPlugin(path.Join(pluginDir, moduleName))
		if err != nil {
			return nil, err
		}
		plugins = append(plugins, pluginInstance)
	}
	return plugins, nil
}

func LoadPlugin(pluginPath string) (*Plugin, error) {
	pluginCfg, err := loadPluginConfig(pluginPath)
	if err != nil {
		return nil, err
	}
	plugin := &Plugin{
		Name:   path.Base(pluginPath),
		Path:   pluginPath,
		Config: pluginCfg,
	}

	if plugin.InitProgram, err = getProgram(pluginPath, "init"); err != nil {
		return nil, fmt.Errorf("Plugin Error (" + pluginPath + " - init: " + err.Error())
	}
	if plugin.DeployProgram, err = getProgram(pluginPath, "deploy"); err != nil {
		return nil, fmt.Errorf("Plugin Error (" + pluginPath + " - deploy: " + err.Error())
	}
	if plugin.SetupProgram, err = getProgram(pluginPath, "setup"); err != nil {
		return nil, fmt.Errorf("Plugin Error (" + pluginPath + " - setup: " + err.Error())
	}
	if plugin.DestroyProgram, err = getProgram(pluginPath, "destroy"); err != nil {
		return nil, fmt.Errorf("Plugin Error (" + pluginPath + " - destroy: " + err.Error())
	}
	return plugin, nil
}

func loadPluginConfig(pluginPath string) (*pluginlib.PluginConfig, error) {
	p := &pluginlib.PluginConfig{}

	yamlFile, err := ioutil.ReadFile(path.Join(pluginPath, "stackhead-module.yml"))
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(yamlFile, &p); err != nil {
		return nil, err
	}

	// Add plugin path to Init path
	if p.Terraform.Provider.Init != "" {
		p.Terraform.Provider.Init = path.Join(pluginPath, p.Terraform.Provider.Init)
		// Ensure new path is still in plugin path
		if !strings.HasPrefix(p.Terraform.Provider.Init, pluginPath) {
			return nil, fmt.Errorf("path security violated: Init path does not resolve to plugin path")
		}
	}

	return p, nil
}

func getProgram(path string, fileName string) (*PluginProgram, error) {
	fullPath := path + "/" + fileName + ".js"
	src, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	program := &PluginProgram{
		OriginalPath: fullPath,
		Source:       string(src),
	}
	return program, err
}
