package container_docker

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	xfs "github.com/saitho/golang-extended-fs/v2"

	"github.com/getstackhead/stackhead/config"
	"github.com/getstackhead/stackhead/system"
)

func InstallApt() error {
	// Add Docker apt signing key
	if _, _, err := system.RemoteRun("mkdir -p /etc/apt/keyrings", system.RemoteRunOpts{Sudo: true}); err != nil {
		return err
	}
	if _, _, err := system.RemoteRun("curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor --yes -o /etc/apt/keyrings/docker.gpg", system.RemoteRunOpts{}); err != nil {
		return err
	}

	// Setup Docker apt repository on Ubuntu
	if _, _, err := system.RemoteRun("echo \"deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable\" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null", system.RemoteRunOpts{}); err != nil {
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
			Name:   "docker-compose-plugin",
			Vendor: system.PackageVendorApt,
		},
		{
			Name:   "containerd.io",
			Vendor: system.PackageVendorApt,
		},
		{
			Name:   "pass", // for encrypted credentials storage
			Vendor: system.PackageVendorApt,
		},
		{
			Name:   "golang-docker-credential-helpers", // for encrypted credentials storage
			Vendor: system.PackageVendorApt,
		},
	}); err != nil {
		return err
	}
	return nil
}

func PrepareFiles() error {
	// Add stackhead user to docker
	if _, _, err := system.RemoteRun("usermod", system.RemoteRunOpts{Args: []string{"-a -G docker stackhead"}}); err != nil {
		return fmt.Errorf("unable to add stackhead user to docker group")
	}

	// Create .docker folder
	if err := xfs.CreateFolder("ssh://" + path.Join(config.RootDirectory, ".docker")); err != nil {
		return err
	}
	if err := xfs.Chown("ssh://"+path.Join(config.RootDirectory, ".docker"), 1412, 1412); err != nil {
		return err
	}
	return nil
}

func lookupCredentialsKey() (string, error) {
	stdout, err := system.SimpleRemoteRun("gpg", system.RemoteRunOpts{User: "stackhead", Args: []string{"-k \"StackHead Docker Credentials Store\""}})
	if err != nil {
		if strings.Contains(err.Error(), "No public key") {
			// no keys found at all which may happen when GPG is used for the first time; do not consider an error
			return "", nil
		}
		return "", err
	}
	var re = regexp.MustCompile(`(?m)^\s+(\w+)\nuid\s+\[ultimate\] StackHead Docker Credentials Store$`)
	matches := re.FindAllStringSubmatch(stdout, -1)
	if len(matches) == 0 {
		return "", nil
	}
	if len(matches) > 1 {
		return "", fmt.Errorf("more than one valid GPG key with name \"StackHead Docker Credentials Store\" found.")
	}
	return matches[0][1], nil
}

func PrepareCredentialsStore() error {
	// Look for existing credentials GPG key
	key, err := lookupCredentialsKey()
	if err != nil {
		return fmt.Errorf("unable to lookup GPG key for Docker Credentials Store: " + err.Error())
	}
	if key == "" {
		// Key does not exist: Create GPG key for credentails encryption
		_, err := system.SimpleRemoteRun("gpg", system.RemoteRunOpts{User: "stackhead", Args: []string{"--batch --gen-key <<EOF\n\tKey-Type: 1\n\tKey-Length: 2048\n\tSubkey-Type: 1\n\tSubkey-Length: 2048\n\tName-Real: StackHead Docker Credentials Store\n\tExpire-Date: 0\n\t%no-protection\nEOF"}})
		if err != nil {
			return fmt.Errorf("unable to create GPG key for Docker Credentials Store: " + err.Error())
		}
		key, err = lookupCredentialsKey()
		if err != nil {
			return fmt.Errorf("unable to lookup GPG key for Docker Credentials Store: " + err.Error())
		}
	}
	fmt.Println("Docker Credentials encryption key is " + key)
	// todo: log created resource

	_, err = system.SimpleRemoteRun("pass", system.RemoteRunOpts{User: "stackhead", Args: []string{"init", key}})
	if err != nil {
		return err
	}

	// todo: consider existing config files on systems not set up via StackHead
	if err := xfs.WriteFile("ssh://"+path.Join(config.RootDirectory, ".docker", "config.json"), "{\"credStore\":\"pass\"}"); err != nil {
		return err
	}
	if err := xfs.Chown("ssh://"+path.Join(config.RootDirectory, ".docker", "config.json"), 1412, 1412); err != nil {
		return err
	}

	return nil
}

func (Module) Install(modulesSettings interface{}) error {
	if err := InstallApt(); err != nil {
		return err
	}
	if err := PrepareFiles(); err != nil {
		return err
	}
	if err := PrepareCredentialsStore(); err != nil {
		return err
	}
	return nil
}
