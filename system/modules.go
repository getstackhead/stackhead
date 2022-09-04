package system

import (
	"bytes"
	"encoding/json"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/hairyhenderson/gomplate/v3"
	xfs "github.com/saitho/golang-extended-fs"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"

	"github.com/getstackhead/stackhead/project"
)

type ModuleConfig struct {
	Name      string
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
	InitFuncMap        template.FuncMap
}

func (a ModuleTerraformConfigProvider) Equal(b ModuleTerraformConfigProvider) bool {
	return (a.Vendor == b.Vendor) && (a.Name == b.Name) && (a.Version == b.Version)
}

type ModuleTerraformConfig struct {
	Provider ModuleTerraformConfigProvider
}

type Module interface {
	Install(moduleSettings interface{}) error
	Deploy(moduleSettings interface{}) error
	GetConfig() ModuleConfig
}

type ModuleTemplateData struct {
	Project *project.Project
}

func RenderModuleTemplate(fileName string, additionalData map[string]interface{}, additionalFuncs template.FuncMap) (string, error) {
	fileContents, err := xfs.ReadFile("pkging:///templates/modules/" + fileName)
	if err != nil {
		return "", err
	}
	return RenderModuleTemplateText(fileName, fileContents, additionalData, additionalFuncs)
}

func RenderModuleTemplateText(templateName string, fileContents string, additionalData map[string]interface{}, additionalFuncs template.FuncMap) (string, error) {
	tmpl := template.New(templateName)

	// prepare functions
	tmpl = tmpl.Funcs(sprig.TxtFuncMap()).Funcs(gomplate.CreateFuncs(nil, nil))
	if additionalFuncs != nil {
		tmpl = tmpl.Funcs(additionalFuncs)
	}

	// prepare data
	data := map[string]interface{}{
		"Project": Context.Project,
	}
	if additionalData != nil {
		maps.Copy(data, additionalData)
	}

	tmpl, err := tmpl.Parse(fileContents)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func GetModuleSettings(moduleName string) interface{} {
	return viper.GetStringMap("modules_config")[moduleName]
}

func UnpackModuleSettings[T interface{}](_modulesSettings interface{}) (*T, error) {
	dbByte, _ := json.Marshal(_modulesSettings.(map[string]interface{}))
	modulesSettings := new(T)
	err := json.Unmarshal(dbByte, &modulesSettings)
	return modulesSettings, err
}
