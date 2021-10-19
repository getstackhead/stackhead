package declarations

import (
	"github.com/getstackhead/stackhead/pluginlib"
	"github.com/getstackhead/stackhead/stackhead"
)

var InstallPackage = func(packages []pluginlib.Package) error {
	for _, p := range packages {
		if p.Vendor == pluginlib.PackageVendorApt {
			if err := stackhead.RemoteRun("apt install " + p.Name); err != nil {
				return err
			}
		} else if p.Vendor == pluginlib.PackageVendorApk {
			if err := stackhead.RemoteRun("apk add " + p.Name); err != nil {
				return err
			}
		}
	}
	return nil
}
