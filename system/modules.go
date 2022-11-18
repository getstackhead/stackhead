package system

import (
	"bytes"
	"embed"
	"encoding/json"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/hairyhenderson/gomplate/v3"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"

	"github.com/getstackhead/stackhead/project"
)

type ModuleConfig struct {
	Name string
	Type string
}

type Module interface {
	Install(moduleSettings interface{}) error
	Deploy(moduleSettings interface{}) error
	Destroy(moduleSettings interface{}) error
	Init(moduleSettings interface{})
	GetConfig() ModuleConfig
	GetTemplates() *embed.FS
}

type PluginModule interface {
	Init(moduleSettings interface{})
	GetConfig() ModuleConfig
	GetTemplates() *embed.FS
}

type ModuleTemplateData struct {
	Project *project.Project
}

func RenderModuleTemplate(templateFolder embed.FS, fileName string, additionalData map[string]any, additionalFuncs template.FuncMap) (string, error) {
	fileContents, err := templateFolder.ReadFile("templates/" + fileName)
	if err != nil {
		return "", err
	}
	return RenderModuleTemplateText(fileName, string(fileContents), additionalData, additionalFuncs)
}

func RenderModuleTemplateText(templateName string, fileContents string, additionalData map[string]any, additionalFuncs template.FuncMap) (string, error) {
	tmpl := template.New(templateName)

	// prepare functions
	tmpl = tmpl.Funcs(sprig.TxtFuncMap()).Funcs(gomplate.CreateFuncs(nil, nil))
	tmpl = tmpl.Funcs(map[string]any{
		"joinMap": func(data map[string]interface{}, joinChar string, wrapChar string) []string {
			var list []string
			for key := range data {
				list = append(list, wrapChar+key+"="+data[key].(string)+wrapChar)
			}
			return list
		},
	})
	if additionalFuncs != nil {
		tmpl = tmpl.Funcs(additionalFuncs)
	}

	// prepare data
	data := map[string]any{
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
