package system

type PackageVendor string

var PackageVendorApt PackageVendor = "apt"
var PackageVendorApk PackageVendor = "apk"

type Package struct {
	Name   string
	Vendor PackageVendor
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
