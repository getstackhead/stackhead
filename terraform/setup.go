package terraform

import (
	"bytes"
	"embed"
	"fmt"
	"path"
	"text/template"

	"github.com/google/go-cmp/cmp"
	"github.com/gookit/event"
	xfs "github.com/saitho/golang-extended-fs/v2"
	logger "github.com/sirupsen/logrus"

	"github.com/getstackhead/stackhead/config"
	"github.com/getstackhead/stackhead/project"
	"github.com/getstackhead/stackhead/system"
)

var terraformProvidersFile = path.Join(config.RootTerraformDirectory, "terraform-providers.tf")

//go:embed templates
var terraformTemplates embed.FS

func Setup() error {
	event.MustFire("setup.terraform.pre-install", event.M{})
	if err := xfs.CreateFolder("ssh://" + config.RootTerraformDirectory); err != nil {
		return err
	}
	if err := xfs.Chown("ssh://"+config.RootTerraformDirectory, 1412, 1412); err != nil {
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
	event.MustFire("setup.terraform.post-install", event.M{})

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
}

func FilterModulesWithProviders(modules []system.Module) []system.Module {
	var modulesWithProviders []system.Module
	emptyProvider := system.ModuleTerraformConfigProvider{}
	for _, module := range modules {
		moduleCfg := module.GetConfig()
		if cmp.Equal(moduleCfg.Terraform.Provider, emptyProvider) {
			continue
		}
		modulesWithProviders = append(modulesWithProviders, module)
	}
	return modulesWithProviders
}

func BuildAndWriteProviders() error {
	modulesWithProviders := FilterModulesWithProviders(system.Context.GetModulesInOrder())
	fileContents, err := BuildProviders(modulesWithProviders, NO_PER_PROJECT)
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
	// full access to stackhead user to Terraform folder
	if _, _, err := system.RemoteRun("chown", "-R", "stackhead:stackhead", path.Join(config.RootTerraformDirectory)); err != nil {
		return err
	}
	// keep root permissions on base file terraform-providers.tf
	if _, _, err := system.RemoteRun("chown", "-R", "root:root", path.Join(config.RootTerraformDirectory, "terraform-providers.tf")); err != nil {
		return err
	}
	return nil
}

type BuildProviderMode int

var ONLY_PER_PROJECT BuildProviderMode = 1
var NO_PER_PROJECT BuildProviderMode = 2

func BuildProviders(modules []system.Module, mode BuildProviderMode) (bytes.Buffer, error) {
	var providers []system.ModuleTerraformConfigProvider
	for _, module := range modules {
		providers = append(providers, module.GetConfig().Terraform.Provider)
	}

	data := Data{
		Providers:         providers,
		Context:           system.Context,
		AdditionalContent: "",
	}

	// Additional provider configuration from plugins
	var suffix string
	for _, module := range modules {
		provider := module.GetConfig().Terraform.Provider
		if (mode == ONLY_PER_PROJECT && !provider.ProviderPerProject) || (mode == NO_PER_PROJECT && provider.ProviderPerProject) {
			continue
		}
		if provider.Init != "" {
			templateContent, err := module.GetTemplates().ReadFile("templates/" + provider.Init)
			if err != nil {
				return bytes.Buffer{}, err
			}
			providerInitContent, err := buildProvider(string(templateContent), data, provider.InitFuncMap)
			if err != nil {
				return bytes.Buffer{}, fmt.Errorf("Unable to load module init file \"" + provider.Init + "\": " + err.Error())
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
	if mode == ONLY_PER_PROJECT {
		returnBuf := bytes.Buffer{}
		returnBuf.WriteString(data.AdditionalContent)
		return returnBuf, nil
	}

	templateContent, err := terraformTemplates.ReadFile("templates/terraform-providers.tf.tmpl")
	if err != nil {
		return bytes.Buffer{}, err
	}
	return buildProvider(string(templateContent), data, nil)
}

func buildProvider(templateContent string, data Data, funcMap template.FuncMap) (bytes.Buffer, error) {
	baseTmpl := template.New("providers")

	if funcMap != nil {
		baseTmpl.Funcs(funcMap)
	}

	var buf bytes.Buffer
	tmpl, err := baseTmpl.Parse(templateContent)
	if err != nil {
		return buf, err
	}
	err = tmpl.Execute(&buf, data)
	return buf, err
}
