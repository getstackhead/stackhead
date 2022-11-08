package container_docker

import (
	"fmt"
	"github.com/getstackhead/stackhead/system"
)

func InstallApt() error {
	// Add Docker apt signing key
	if _, _, err := system.RemoteRun("sudo mkdir -p /etc/apt/keyrings"); err != nil {
		return err
	}
	if _, _, err := system.RemoteRun("curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor --yes -o /etc/apt/keyrings/docker.gpg"); err != nil {
		return err
	}

	// Setup Docker apt repository on Ubuntu
	if _, _, err := system.RemoteRun("echo \"deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable\" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null"); err != nil {
		return err
	}
	if err := system.UpdatePackageList(system.PackageVendorApt); err != nil {
		return err
	}

	// Install Docker-CE
	if err := system.InstallPackage([]system.Package{
		{
			Name:   "docker-ce",
			Vendor: system.PackageVendorApt,
		},
		{
			Name:   "docker-ce-cli",
			Vendor: system.PackageVendorApt,
		},
		{
			Name:   "containerd.io",
			Vendor: system.PackageVendorApt,
		},
	}); err != nil {
		return err
	}

	// Add stackhead user to docker
	if _, _, err := system.RemoteRun("usermod", "-a -G docker stackhead"); err != nil {
		return fmt.Errorf("unable to add stackhead user to docker group")
	}

	return nil
}

func (Module) Install(modulesSettings interface{}) error {
	return InstallApt()
}
