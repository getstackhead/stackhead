package plugins

import (
	"fmt"
	"io"
	"reflect"

	"github.com/open2b/scriggo/native"

	"github.com/getstackhead/stackhead/pluginlib"
	"github.com/getstackhead/stackhead/plugins/declarations"
)

func getPluginLib() native.Package {
	return native.Package{
		Name: "pluginlib",
		Declarations: native.Declarations{
			"StackHeadExecute": declarations.StackHeadExecute,
			"InstallPackage":   declarations.InstallPackage,
			"GetProject":       declarations.GetProject,

			// Structs
			"Package":          reflect.TypeOf(pluginlib.Package{}),
			"PackageVendorApt": pluginlib.PackageVendorApt,
			"PackageVendorApk": pluginlib.PackageVendorApk,
		},
	}
}

func getPackages() native.Packages {
	return native.Packages{
		"github.com/getstackhead/pluginlib": getPluginLib(),
		"fmt": native.Package{
			Name: "fmt",
			Declarations: native.Declarations{
				"Print": func(a ...interface{}) (int, error) {
					return 0, nil
					//return fmt.Print(append([]interface{}{"[StackHead Plugin]"}, a...)...)
				},
				"Println": func(a ...interface{}) (int, error) {
					return 0, nil
					//return fmt.Println(append([]interface{}{"[StackHead Plugin]"}, a...)...)
				},
				"Printf": func(format string, a ...interface{}) (int, error) {
					return 0, nil
					//return fmt.Printf("[StackHead Plugin] " + format, a...)
				},
				"Errorf": func(format string, a ...interface{}) error {
					return nil
					//return fmt.Errorf("[StackHead Plugin] " + format, a...)
				},
				"Fprintf": func(w io.Writer, format string, a ...interface{}) (n int, err error) {
					return 0, nil
					//return fmt.Fprintf(w, "[StackHead Plugin] " + format, a...)
				},
				"Fprint": func(w io.Writer, a ...interface{}) (n int, err error) {
					return 0, nil
					//return fmt.Fprint(w, append([]interface{}{"[StackHead Plugin]"}, a...)...)
				},
				"Fprintln": func(w io.Writer, a ...interface{}) (n int, err error) {
					return 0, nil
					//return fmt.Fprintln(w, append([]interface{}{"[StackHead Plugin]"}, a...)...)
				},
				"Sprintf":  fmt.Sprintf,
				"Sprint":   fmt.Sprint,
				"Sprintln": fmt.Sprintf,
				"Fscan":    fmt.Fscan,
				"Fscanln":  fmt.Fscanln,
				"Fscanf":   fmt.Fscanf,
				"Sscan":    fmt.Sscan,
				"Sscanln":  fmt.Sscanln,
				"Sscanf":   fmt.Sscanf,
				"Scan":     fmt.Scan,
				"Scanln":   fmt.Scanln,
				"Scanf":    fmt.Scanf,
			},
		},
	}
}
