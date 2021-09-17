package ansible

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"

	"github.com/getstackhead/stackhead/stackhead"
)

type Inventory struct {
	All struct {
		Vars struct {
			AnsibleUser               string                 `yaml:"ansible_user"`
			AnsibleConnection         string                 `yaml:"ansible_connection"`
			AnsiblePythonInterpreter  string                 `yaml:"ansible_python_interpreter"`
			StackHeadConfigFolder     string                 `yaml:"stackhead__config_folder"`
			StackHeadDNS              []string               `yaml:"stackhead__dns"`
			StackHeadWebserver        string                 `yaml:"stackhead__webserver"`
			StackHeadContainer        string                 `yaml:"stackhead__container"`
			StackHeadPlugins          []string               `yaml:"stackhead__plugins"`
			StackHeadConfigSetup      map[string]interface{} `yaml:"stackhead__config_setup"`
			StackHeadConfigDeployment map[string]interface{} `yaml:"stackhead__config_deployment"`
			StackHeadConfigDestroy    map[string]interface{} `yaml:"stackhead__config_destroy"`
			CertificatesEmailAddress  string                 `yaml:"certificates_email_address"`
			TerraformUpdateInterval   string                 `yaml:"stackhead__tf_update_interval"`
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

// CreateInventoryFile creates a temporary inventory file with the given settings and returns an absolute file path.
// Note: make sure to remove the file when you don't need it anymore!
func CreateInventoryFile(ipAddress string, projectDefinitionFile string) (string, error) {
	var err error

	conf := Inventory{}
	conf.All.Vars.AnsibleUser = "root"
	conf.All.Vars.AnsibleConnection = "ssh"
	conf.All.Vars.AnsiblePythonInterpreter = "/usr/bin/python3"
	conf.All.Hosts.Mackerel.AnsibleHost = ipAddress

	conf.All.Vars.StackHeadWebserver, err = stackhead.GetWebserverModule()
	if err != nil {
		return "", err
	}

	conf.All.Vars.StackHeadContainer, err = stackhead.GetContainerModule()
	if err != nil {
		return "", err
	}

	conf.All.Vars.StackHeadDNS, err = stackhead.GetDNSModules()
	if err != nil {
		return "", err
	}

	conf.All.Vars.StackHeadPlugins, err = stackhead.GetPluginModules()
	if err != nil {
		return "", err
	}

	tmpFile, err := ioutil.TempFile("", "inventory")
	if err != nil {
		return "", err
	}

	inventoryFilePath, err := filepath.Abs(tmpFile.Name())
	if err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", err
	}

	if len(projectDefinitionFile) > 0 {
		conf.All.Vars.StackHeadConfigFolder, err = filepath.Abs(filepath.Dir(projectDefinitionFile))
		if err != nil {
			return "", err
		}
		projectDefinitionFilePath := filepath.Base(projectDefinitionFile)

		projectDefinitionFilePath = strings.TrimSuffix(projectDefinitionFilePath, ".stackhead.yml")
		projectDefinitionFilePath = strings.TrimSuffix(projectDefinitionFilePath, ".stackhead.yaml")
		conf.All.Hosts.Mackerel.Stackhead.Applications = append(conf.All.Hosts.Mackerel.Stackhead.Applications, projectDefinitionFilePath)
	}

	conf.All.Vars.StackHeadConfigSetup = viper.GetStringMap("config.setup")
	conf.All.Vars.StackHeadConfigDeployment = viper.GetStringMap("config.deployment")
	conf.All.Vars.StackHeadConfigDestroy = viper.GetStringMap("config.destroy")

	if viper.IsSet("terraform.update_interval") {
		conf.All.Vars.TerraformUpdateInterval = viper.GetString("terraform.update_interval")
	}

	if viper.IsSet("certificates.register_email") {
		conf.All.Vars.CertificatesEmailAddress = viper.GetString("certificates.register_email")
	}

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

	return inventoryFilePath, nil
}
