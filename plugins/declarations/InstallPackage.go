package declarations

import (
	"github.com/getstackhead/stackhead/pluginlib"
	"github.com/getstackhead/stackhead/system"
)

var InstallPackage = func(packages []pluginlib.Package) error {
	for _, p := range packages {
		if p.Vendor == pluginlib.PackageVendorApt {
			if _, _, err := system.RemoteRun("DEBIAN_FRONTEND=noninteractive apt", "install -yq "+p.Name); err != nil {
				return err
			}
		} else if p.Vendor == pluginlib.PackageVendorApk {
			if _, _, err := system.RemoteRun("apk", "add "+p.Name); err != nil {
				return err
			}
		}
	}
	return nil
}
