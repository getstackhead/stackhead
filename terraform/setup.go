package terraform

import (
	"bytes"
	xfs "github.com/saitho/golang-extended-fs"
	"path"
	"text/template"

	logger "github.com/sirupsen/logrus"

	"github.com/getstackhead/stackhead/config"
	"github.com/getstackhead/stackhead/pluginlib"
	"github.com/getstackhead/stackhead/plugins"
	"github.com/getstackhead/stackhead/plugins/declarations"
	"github.com/getstackhead/stackhead/system"
)

var terraformProvidersFile = path.Join(config.Paths.RootTerraform, "terraform-providers.tf")

func Setup() error {
	if err := xfs.CreateFolder("ssh://" + config.Paths.RootTerraform); err != nil {
		return err
	}

	if _, _, err := declarations.StackHeadExecute("curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -"); err != nil {
		return err
	}
	if _, _, err := declarations.StackHeadExecute("sudo apt-add-repository \"deb [arch=$(dpkg --print-architecture)] https://apt.releases.hashicorp.com $(lsb_release -cs) main\""); err != nil {
		return err
	}
	if err := declarations.InstallPackage([]pluginlib.Package{
		{
			Name:   "terraform=1.0.9",
			Vendor: "apt",
		},
	}); err != nil {
		return err
	}

	//# Setup Terraform components of modules
	//- import_tasks: "../roles/stackhead_module_api/tasks_internal/terraform.yml"
	//- import_tasks: "../roles/config_terraform/tasks/generate-providers.yml"
	//
	//# Initial run applying settings from Terraform files created above
	//- import_tasks: "../roles/config_terraform/tasks/execute.yml"
	return nil
}

type Data struct {
	Providers []pluginlib.PluginTerraformConfigProvider
	Context   system.ContextStruct

	AdditionalContent string

	// SSL certificate information
	SnakeoilFullchainPath string
	SnakeoilPrivkeyPath   string
}

func BuildAndWriteProviders(p []*plugins.Plugin) error {
	var providers []pluginlib.PluginTerraformConfigProvider
	for _, plugin := range p {
		emptyProvider := pluginlib.PluginTerraformConfigProvider{}
		if plugin.Config.Terraform.Provider == emptyProvider {
			continue
		}
		providers = append(providers, plugin.Config.Terraform.Provider)
	}
	fileContents, err := buildProviders(providers)
	if err != nil {
		return err
	}
	// Write file
	if err := xfs.WriteFile("ssh://"+terraformProvidersFile, fileContents.String()); err != nil {
		return err
	}

	// todo: symlink provider configuration into project directories

	return nil
}

func SymlinkProviders(project *pluginlib.Project) error {
	_, errMsg, err := system.RemoteRun("ln", "-sf "+terraformProvidersFile+" "+path.Join(config.Paths.GetProjectTerraformDirectoryPath(project), "terraform-providers.tf"))
	if err != nil {
		logger.Errorln(errMsg.String())
	}
	return err
}

func Init(directory string) error {
	if _, _, err := system.RemoteRun("(cd " + directory + " && " + GetCommand("init") + ")"); err != nil {
		return err
	}
	return nil
}

func Apply(directory string) error {
	if _, _, err := system.RemoteRun("(cd " + directory + " && " + GetCommand("apply -auto-approve") + ")"); err != nil {
		return err
	}
	return nil
}

func InstallProviders() error {
	if err := Init(config.Paths.RootTerraform); err != nil {
		return err
	}
	if err := Apply(config.Paths.RootTerraform); err != nil {
		return err
	}
	if err := xfs.Chown("ssh://"+config.Paths.RootTerraform, 1412, 1412); err != nil {
		return err
	}
	SnakeoilFullchainPath := config.Paths.GetSnakeoilFullchainPath()
	SnakeoilPrivkeyPath := config.Paths.GetSnakeoilPrivKeyPath()
	if err := xfs.Chown("ssh://"+SnakeoilFullchainPath, 1412, 1412); err != nil {
		return err
	}
	if err := xfs.Chown("ssh://"+SnakeoilPrivkeyPath, 1412, 1412); err != nil {
		return err
	}
	return nil
}

func buildProviders(providers []pluginlib.PluginTerraformConfigProvider) (bytes.Buffer, error) {
	data := Data{
		Providers:             providers,
		Context:               system.Context,
		SnakeoilFullchainPath: config.Paths.GetSnakeoilFullchainPath(),
		SnakeoilPrivkeyPath:   config.Paths.GetSnakeoilPrivKeyPath(),
		AdditionalContent:     "",
	}

	// Additional provider configuration from plugins
	var suffix string
	for _, provider := range providers {
		if provider.ProviderPerProject {
			continue
		}
		if provider.Init != "" {
			// todo: load template from .Init
			providerInitContent, err := buildProvider(provider.Init, data)
			if err != nil {
				return bytes.Buffer{}, err
			}
			data.AdditionalContent += "\n" + providerInitContent.String() + "\n"
		} else {
			suffix = ""
			if provider.NameSuffix != "" {
				suffix = "-" + provider.NameSuffix
			}
			data.AdditionalContent += "\nprovider \"" + provider.Name + suffix + "\" {\n}\n"
		}
	}
	return buildProvider("pkging:///templates/terraform-providers.tf.tmpl", data)
}

func buildProvider(filePath string, data Data) (bytes.Buffer, error) {
	var buf bytes.Buffer
	fileContents, err := xfs.ReadFile(filePath)
	if err != nil {
		return buf, err
	}
	tmpl, err := template.New("providers").Parse(fileContents)
	if err != nil {
		return buf, err
	}
	err = tmpl.Execute(&buf, data)
	return buf, err
}
