package ansible

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/getstackhead/stackhead/cli/stackhead"
)

type Inventory struct {
	All struct {
		Vars struct {
			AnsibleUser               string                 `yaml:"ansible_user"`
			AnsibleConnection         string                 `yaml:"ansible_connection"`
			AnsiblePythonInterpreter  string                 `yaml:"ansible_python_interpreter"`
			StackHeadConfigFolder     string                 `yaml:"stackhead__config_folder"`
			StackHeadWebserver        string                 `yaml:"stackhead__webserver"`
			StackHeadContainer        string                 `yaml:"stackhead__container"`
			StackHeadPlugins          []string               `yaml:"stackhead__plugins"`
			StackHeadConfigSetup      map[string]interface{} `yaml:"stackhead__config_setup"`
			StackHeadConfigDeployment map[string]interface{} `yaml:"stackhead__config_deployment"`
			StackHeadConfigDestroy    map[string]interface{} `yaml:"stackhead__config_destroy"`
		}
		Hosts struct {
			Mackerel struct {
				AnsibleHost string `yaml:"ansible_host"`
				Stackhead   struct {
					Applications []string
				}
			}
		}
	}
}

func copyFile(srcPath string, destPath string) error {
	input, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(destPath, input, 0600); err != nil {
		return err
	}
	return nil
}

// CreateInventoryFile creates a temporary inventory file with the given settings and returns an absolute file path.
// Note: make sure to remove the file when you don't need it anymore!
func CreateInventoryFile(ipAddress string, projectDefinitionFile string) (string, error) {
	conf := Inventory{}
	conf.All.Vars.AnsibleUser = "root"
	conf.All.Vars.AnsibleConnection = "ssh"
	conf.All.Vars.AnsiblePythonInterpreter = "/usr/bin/python3"
	conf.All.Hosts.Mackerel.AnsibleHost = ipAddress

	webserverModule, err := stackhead.GetWebserverModule()
	if err != nil {
		return "", err
	}
	conf.All.Vars.StackHeadWebserver = webserverModule

	containerModule, err := stackhead.GetContainerModule()
	if err != nil {
		return "", err
	}
	conf.All.Vars.StackHeadContainer = containerModule

	pluginModules, err := stackhead.GetPluginModules()
	if err != nil {
		return "", err
	}
	conf.All.Vars.StackHeadPlugins = pluginModules

	tmpFile, err := ioutil.TempFile("", "inventory")
	if err != nil {
		return "", err
	}

	filePath, err := filepath.Abs(tmpFile.Name())
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	if len(projectDefinitionFile) > 0 {
		conf.All.Vars.StackHeadConfigFolder, err = ioutil.TempDir("", "project_definitions")
		if err != nil {
			return "", err
		}
		projectDefinitionFilePath := filepath.Base(projectDefinitionFile)

		// Copy project definition file
		err = copyFile(projectDefinitionFile, filepath.Join(conf.All.Vars.StackHeadConfigFolder, projectDefinitionFilePath))
		if err != nil {
			return "", err
		}

		projectDefinitionFilePath = strings.TrimSuffix(projectDefinitionFilePath, ".yml")
		projectDefinitionFilePath = strings.TrimSuffix(projectDefinitionFilePath, ".yaml")
		conf.All.Hosts.Mackerel.Stackhead.Applications = append(conf.All.Hosts.Mackerel.Stackhead.Applications, projectDefinitionFilePath)
	}

	conf.All.Vars.StackHeadConfigSetup = viper.GetStringMap("config.setup")
	conf.All.Vars.StackHeadConfigDeployment = viper.GetStringMap("config.deployment")
	conf.All.Vars.StackHeadConfigDestroy = viper.GetStringMap("config.destroy")

	d, err := yaml.Marshal(&conf)
	if err != nil {
		return "", err
	}

	if _, err = tmpFile.Write(d); err != nil {
		return "", err
	}

	// Close the file
	if err = tmpFile.Close(); err != nil {
		return "", err
	}

	return filePath, nil
}
