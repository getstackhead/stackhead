package terraform

import (
	"bytes"
	"path"
	"text/template"

	xfs "github.com/saitho/golang-extended-fs"
	logger "github.com/sirupsen/logrus"

	"github.com/getstackhead/stackhead/config"
	"github.com/getstackhead/stackhead/project"
	"github.com/getstackhead/stackhead/system"
)

var terraformProvidersFile = path.Join(config.RootTerraformDirectory, "terraform-providers.tf")

func Setup() error {
	if err := xfs.CreateFolder("ssh://" + config.RootTerraformDirectory); err != nil {
		return err
	}

	if _, _, err := system.RemoteRun("curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -"); err != nil {
		return err
	}
	if _, _, err := system.RemoteRun("sudo apt-add-repository \"deb [arch=$(dpkg --print-architecture)] https://apt.releases.hashicorp.com $(lsb_release -cs) main\""); err != nil {
		return err
	}
	if err := system.InstallPackage([]system.Package{
		{
			Name:   "terraform=1.0.9",
			Vendor: system.PackageVendorApt,
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
	Providers []system.ModuleTerraformConfigProvider
	Context   system.ContextStruct

	AdditionalContent string

	// SSL certificate information
	SnakeoilFullchainPath string
	SnakeoilPrivkeyPath   string
}

func BuildAndWriteProviders() error {
	var providers []system.ModuleTerraformConfigProvider
	emptyProvider := system.ModuleTerraformConfigProvider{}
	for _, module := range system.Context.GetModulesInOrder() {
		moduleCfg := module.GetConfig()
		if moduleCfg.Terraform.Provider == emptyProvider {
			continue
		}
		providers = append(providers, moduleCfg.Terraform.Provider)
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

func SymlinkProviders(project *project.Project) error {
	_, errMsg, err := system.RemoteRun("ln", "-sf "+terraformProvidersFile+" "+path.Join(project.GetTerraformDirectoryPath(), "terraform-providers.tf"))
	if err != nil {
		logger.Errorln(errMsg.String())
	}
	return err
}

func Init(directory string) error {
	if _, outErr, err := system.RemoteRun("(cd " + directory + " && " + GetCommand("init") + ")"); err != nil {
		logger.Errorln(outErr.String())
		return err
	}
	return nil
}

func Apply(directory string) error {
	if _, outErr, err := system.RemoteRun("(cd " + directory + " && " + GetCommand("apply -auto-approve") + ")"); err != nil {
		logger.Errorln(outErr.String())
		return err
	}
	return nil
}

func InstallProviders() error {
	if err := Init(config.RootTerraformDirectory); err != nil {
		return err
	}
	if err := Apply(config.RootTerraformDirectory); err != nil {
		return err
	}
	if err := xfs.Chown("ssh://"+config.RootTerraformDirectory, 1412, 1412); err != nil {
		return err
	}
	SnakeoilFullchainPath, SnakeoilPrivkeyPath := config.GetSnakeoilPaths()
	if err := xfs.Chown("ssh://"+SnakeoilFullchainPath, 1412, 1412); err != nil {
		return err
	}
	if err := xfs.Chown("ssh://"+SnakeoilPrivkeyPath, 1412, 1412); err != nil {
		return err
	}
	return nil
}

func buildProviders(providers []system.ModuleTerraformConfigProvider) (bytes.Buffer, error) {
	SnakeoilFullchainPath, SnakeoilPrivkeyPath := config.GetSnakeoilPaths()
	data := Data{
		Providers:             providers,
		Context:               system.Context,
		SnakeoilFullchainPath: SnakeoilFullchainPath,
		SnakeoilPrivkeyPath:   SnakeoilPrivkeyPath,
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
