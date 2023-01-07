package system

import (
	"fmt"
	"strings"
)

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
	groupedPackages := map[PackageVendor][]string{}
	for _, p := range packages {
		groupedPackages[p.Vendor] = append(groupedPackages[p.Vendor], p.Name)
	}

	for vendor, packageNames := range groupedPackages {
		if vendor == PackageVendorApt {
			if _, _, err := RemoteRun("DEBIAN_FRONTEND=noninteractive apt", RemoteRunOpts{Args: []string{"install -yq " + strings.Join(packageNames, " ")}}); err != nil {
				return err
			}
		}
		if vendor == PackageVendorApk {
			if _, _, err := RemoteRun("apk", RemoteRunOpts{Args: []string{"add " + strings.Join(packageNames, " ")}}); err != nil {
				return err
			}
		}
	}

	return nil
}
