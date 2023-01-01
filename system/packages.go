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
			if _, _, err := RemoteRun("apt-get", RemoteRunOpts{Args: []string{"update"}}); err != nil {
				return err
			}
		} else if v == PackageVendorApk {
			if _, _, err := RemoteRun("apk", RemoteRunOpts{Args: []string{"update"}}); err != nil {
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
			if _, _, err := RemoteRun("DEBIAN_FRONTEND=noninteractive apt", RemoteRunOpts{Args: []string{"install -yq " + p.Name}}); err != nil {
				return err
			}
		} else if p.Vendor == PackageVendorApk {
			if _, _, err := RemoteRun("apk", RemoteRunOpts{Args: []string{"add " + p.Name}}); err != nil {
				return err
			}
		}
	}
	return nil
}
