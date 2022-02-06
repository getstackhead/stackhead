package system

type ModuleConfig struct {
	Type      string
	Terraform ModuleTerraformConfig
}

type ModuleTerraformConfigProvider struct {
	Vendor             string
	Name               string
	NameSuffix         string
	Version            string
	ResourceName       string
	ProviderPerProject bool
	Init               string
}

type ModuleTerraformConfig struct {
	Provider ModuleTerraformConfigProvider
}

type Module interface {
	Install() error
	GetConfig() ModuleConfig
}
