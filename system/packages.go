package system

import "fmt"

type PackageVendor string

var PackageVendorApt PackageVendor = "apt"
var PackageVendorApk PackageVendor = "apk"

type Package struct {
	Name   string
	Vendor PackageVendor
}

func UpdatePackageList(vendors ...PackageVendor) error {
	for _, v := range vendors {
		if v == PackageVendorApt {
			if _, _, err := RemoteRun("apt-get", "update"); err != nil {
				return err
			}
		} else if v == PackageVendorApk {
			if _, _, err := RemoteRun("apk", "update"); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unsupported package vendor")
		}
	}
	return nil
}

func InstallPackage(packages []Package) error {
	for _, p := range packages {
		if p.Vendor == PackageVendorApt {
			if _, _, err := RemoteRun("DEBIAN_FRONTEND=noninteractive apt", "install -yq "+p.Name); err != nil {
				return err
			}
		} else if p.Vendor == PackageVendorApk {
			if _, _, err := RemoteRun("apk", "add "+p.Name); err != nil {
				return err
			}
		}
	}
	return nil
}
