package plugins

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/open2b/scriggo"

	"github.com/getstackhead/stackhead/config"
	"github.com/getstackhead/stackhead/pluginlib"
	"github.com/getstackhead/stackhead/stackhead"
)

type Plugin struct {
	Name           string
	Path           string
	Config         *pluginlib.PluginConfig
	InitProgram    *scriggo.Program
	SetupProgram   *scriggo.Program
	DeployProgram  *scriggo.Program
	DestroyProgram *scriggo.Program
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
	return &Plugin{
		Name:           path.Base(pluginPath),
		Path:           pluginPath,
		Config:         pluginCfg,
		InitProgram:    getProgram(pluginPath, "init"),
		DeployProgram:  getProgram(pluginPath, "deploy"),
		SetupProgram:   getProgram(pluginPath, "setup"),
		DestroyProgram: getProgram(pluginPath, "destroy"),
	}, nil
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

func getProgram(path string, fileName string) *scriggo.Program {
	src, err := os.ReadFile(path + "/" + fileName + ".go")
	if err != nil {
		return nil
	}

	// Adapt method name to main()
	var re = regexp.MustCompile(`(?m)^func\s+` + regexp.QuoteMeta(fileName) + `\s*\(\s*\)\s+{$`)
	var count = 1 // negative counter is equivalent to global case (replace all)
	src = []byte(re.ReplaceAllStringFunc(string(src), func(s string) string {
		if count == 0 {
			return s
		}

		count -= 1
		return re.ReplaceAllString(s, "func main() {")
	}))

	fsys := scriggo.Files{"main.go": src}
	opts := &scriggo.BuildOptions{Packages: getPackages()}
	program, err := scriggo.Build(fsys, opts)
	if err != nil {
		fmt.Println("Unable to execute StackHead plugin (" + path + "): " + err.Error())
		return nil
	}
	return program
}
