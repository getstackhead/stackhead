package stackhead_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/getstackhead/stackhead/cli/stackhead"
)

func TestSplitModuleName(t *testing.T) {
	Convey("test splitting module names", t, func() {
		Convey("vendor, type and basename", func() {
			vendor, moduleType, baseName := stackhead.SplitModuleName("getstackhead.stackhead_webserver_nginx")
			So(vendor, ShouldEqual, "getstackhead")
			So(moduleType, ShouldEqual, "stackhead_webserver")
			So(baseName, ShouldEqual, "nginx")
		})
		Convey("only vendor and basename", func() {
			vendor, moduleType, baseName := stackhead.SplitModuleName("getstackhead.nginx")
			So(vendor, ShouldEqual, "getstackhead")
			So(moduleType, ShouldEqual, "")
			So(baseName, ShouldEqual, "nginx")
		})
		Convey("only type and name", func() {
			vendor, moduleType, baseName := stackhead.SplitModuleName("stackhead_webserver_nginx")
			So(vendor, ShouldEqual, "")
			So(moduleType, ShouldEqual, "stackhead_webserver")
			So(baseName, ShouldEqual, "nginx")
		})
		Convey("only basename", func() {
			vendor, moduleType, baseName := stackhead.SplitModuleName("nginx")
			So(vendor, ShouldEqual, "")
			So(moduleType, ShouldEqual, "")
			So(baseName, ShouldEqual, "nginx")
		})
	})
}

func TestGetModuleBaseName(t *testing.T) {
	Convey("test getting module base name", t, func() {
		Convey("vendor, type and basename", func() {
			baseName := stackhead.GetModuleBaseName("getstackhead.stackhead_webserver_nginx")
			So(baseName, ShouldEqual, "nginx")
		})
		Convey("only vendor and basename", func() {
			baseName := stackhead.GetModuleBaseName("getstackhead.nginx")
			So(baseName, ShouldEqual, "nginx")
		})
		Convey("only type and name", func() {
			baseName := stackhead.GetModuleBaseName("stackhead_webserver_nginx")
			So(baseName, ShouldEqual, "nginx")
		})
		Convey("only basename", func() {
			baseName := stackhead.GetModuleBaseName("nginx")
			So(baseName, ShouldEqual, "nginx")
		})
	})
}

func TestExtractVendor(t *testing.T) {
	Convey("test extracting the vendor", t, func() {
		Convey("vendor, type and basename", func() {
			vendor := stackhead.ExtractVendor("getstackhead.stackhead_webserver_nginx")
			So(vendor, ShouldEqual, "getstackhead")
		})
		Convey("only vendor and basename", func() {
			vendor := stackhead.ExtractVendor("getstackhead.nginx")
			So(vendor, ShouldEqual, "getstackhead")
		})
		Convey("no vendor", func() {
			vendor := stackhead.ExtractVendor("stackhead_webserver_nginx")
			So(vendor, ShouldEqual, "")
		})
	})
}

func TestRemoveVendor(t *testing.T) {
	Convey("test removing the vendor", t, func() {
		Convey("vendor, type and basename", func() {
			newName := stackhead.RemoveVendor("getstackhead.stackhead_webserver_nginx")
			So(newName, ShouldEqual, "stackhead_webserver_nginx")
		})
		Convey("only vendor and basename", func() {
			newName := stackhead.RemoveVendor("getstackhead.nginx")
			So(newName, ShouldEqual, "nginx")
		})
		Convey("no vendor", func() {
			newName := stackhead.RemoveVendor("stackhead_webserver_nginx")
			So(newName, ShouldEqual, "stackhead_webserver_nginx")
		})
	})
}

func TestIsWebserverModule(t *testing.T) {
	Convey("test if name is webserver module", t, func() {
		Convey("actual webserver module", func() {
			result := stackhead.IsWebserverModule("getstackhead.stackhead_webserver_nginx")
			So(result, ShouldBeTrue)
		})
		Convey("actual webserver module without vendor", func() {
			result := stackhead.IsWebserverModule("stackhead_webserver_nginx")
			So(result, ShouldBeTrue)
		})
		Convey("no webserver module", func() {
			result := stackhead.IsWebserverModule("getstackhead.stackhead_container_docker")
			So(result, ShouldBeFalse)
		})
		Convey("no webserver module without vendor", func() {
			result := stackhead.IsWebserverModule("stackhead_container_docker")
			So(result, ShouldBeFalse)
		})
		Convey("no module type", func() {
			result := stackhead.IsWebserverModule("docker")
			So(result, ShouldBeFalse)
		})
	})
}

func TestIsContainerModule(t *testing.T) {
	Convey("test if name is container module", t, func() {
		Convey("actual container module", func() {
			result := stackhead.IsContainerModule("getstackhead.stackhead_container_docker")
			So(result, ShouldBeTrue)
		})
		Convey("actual container module without vendor", func() {
			result := stackhead.IsContainerModule("stackhead_container_docker")
			So(result, ShouldBeTrue)
		})
		Convey("no container module", func() {
			result := stackhead.IsContainerModule("getstackhead.stackhead_webserver_nginx")
			So(result, ShouldBeFalse)
		})
		Convey("no container module without vendor", func() {
			result := stackhead.IsContainerModule("stackhead_webserver_nginx")
			So(result, ShouldBeFalse)
		})
		Convey("no module type", func() {
			result := stackhead.IsContainerModule("docker")
			So(result, ShouldBeFalse)
		})
	})
}

func TestGetModuleType(t *testing.T) {
	Convey("get module type of module name", t, func() {
		Convey("webserver module", func() {
			moduleType := stackhead.GetModuleType("getstackhead.stackhead_webserver_nginx")
			So(moduleType, ShouldEqual, stackhead.ModuleWebserver)
		})
		Convey("container module", func() {
			moduleType := stackhead.GetModuleType("getstackhead.stackhead_container_docker")
			So(moduleType, ShouldEqual, stackhead.ModuleContainer)
		})
		Convey("webserver module without vendor", func() {
			moduleType := stackhead.GetModuleType("stackhead_webserver_nginx")
			So(moduleType, ShouldEqual, stackhead.ModuleWebserver)
		})
		Convey("container module without vendor", func() {
			moduleType := stackhead.GetModuleType("stackhead_container_docker")
			So(moduleType, ShouldEqual, stackhead.ModuleContainer)
		})
		Convey("unknown module", func() {
			moduleType := stackhead.GetModuleType("getstackhead.stackhead_unknown_docker")
			So(moduleType, ShouldEqual, "")
		})
	})
}

func TestAutoCompleteModuleName(t *testing.T) {
	Convey("complete module name", t, func() {
		Convey("already complete module", func() {
			oldName := "getstackhead.stackhead_webserver_nginx"
			name, err := stackhead.AutoCompleteModuleName(oldName, stackhead.ModuleWebserver)
			So(name, ShouldEqual, oldName)
			So(err, ShouldBeNil)
		})
		Convey("completed module does not match expected type", func() {
			oldName := "getstackhead.stackhead_webserver_nginx"
			name, err := stackhead.AutoCompleteModuleName(oldName, stackhead.ModuleContainer)
			So(name, ShouldBeEmpty)
			So(err, ShouldNotBeNil)
		})
		Convey("prepend default vendor", func() {
			oldName := "stackhead_webserver_nginx"
			name, err := stackhead.AutoCompleteModuleName(oldName, stackhead.ModuleWebserver)
			So(name, ShouldEqual, "getstackhead.stackhead_webserver_nginx")
			So(err, ShouldBeNil)
		})
		Convey("prepend type", func() {
			oldName := "randomvendor.nginx"
			name, err := stackhead.AutoCompleteModuleName(oldName, stackhead.ModuleWebserver)
			So(name, ShouldEqual, "randomvendor.stackhead_webserver_nginx")
			So(err, ShouldBeNil)
		})
		Convey("prepend default vendor and type", func() {
			oldName := "nginx"
			name, err := stackhead.AutoCompleteModuleName(oldName, stackhead.ModuleWebserver)
			So(name, ShouldEqual, "getstackhead.stackhead_webserver_nginx")
			So(err, ShouldBeNil)
		})
	})
}
