package plugins

import (
	"reflect"

	"github.com/getstackhead/stackhead/config"
	"github.com/getstackhead/stackhead/pluginlib"
	"github.com/getstackhead/stackhead/plugins/declarations"
)

type StackheadPluginFuncsType map[string]interface{}

type PathsStruct struct {
	RootDirectory          string
	CertificatesDirectory  string
	RootTerraformDirectory string
	ProjectsRootDirectory  string
}

var StackheadPluginFuncs = StackheadPluginFuncsType{
	"StackHeadExecute":    declarations.StackHeadExecute,
	"InstallPackage":      declarations.InstallPackage,
	"GetProject":          declarations.GetProject,
	"RenderTemplate":      RenderTemplate,
	"CreateTerraformFile": CreateTerraformFile,

	// Structs
	"Package": reflect.TypeOf(pluginlib.Package{}),
	"Paths":   config.Paths,
}
