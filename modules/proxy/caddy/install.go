package proxy_caddy

import (
	"fmt"
	xfs "github.com/saitho/golang-extended-fs"
	logger "github.com/sirupsen/logrus"

	"github.com/getstackhead/stackhead/system"
)

func InstallApt() error {
	// Install packages to allow apt to use a repository over HTTPS
	if err := system.InstallPackage([]system.Package{
		{Name: "apt-transport-https", Vendor: system.PackageVendorApt},
		{Name: "debian-keyring", Vendor: system.PackageVendorApt},
		{Name: "debian-archive-keyring", Vendor: system.PackageVendorApt},
	}); err != nil {
		return err
	}

	// Add Caddy apt signing key
	if _, _, err := system.RemoteRun("curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | \nWorkaround; please use this until the docs update is pushed live:\n\ncurl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg"); err != nil {
		return err
	}

	// Setup Caddy apt repository on Ubuntu
	if _, _, err := system.RemoteRun("curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list"); err != nil {
		return err
	}
	if err := system.UpdatePackageList(system.PackageVendorApt); err != nil {
		return err
	}

	// Install Caddy
	if err := system.InstallPackage([]system.Package{
		{
			Name:   "caddy",
			Vendor: system.PackageVendorApt,
		},
	}); err != nil {
		return err
	}

	// Create /etc/caddy/conf.d/ folder
	if err := xfs.CreateFolder("ssh:///etc/caddy/conf.d"); err != nil {
		return fmt.Errorf("unable to create Caddy conf.d folder: " + err.Error())
	}
	if err := xfs.Chown("ssh:///etc/caddy/conf.d", 1412, 1412); err != nil {
		return fmt.Errorf("unable to change owner of Caddy conf.d folder: " + err.Error())
	}

	// Overwrite Caddyfile

	// todo: make configurable and supply caddy config as data
	caddyFile, err := system.RenderModuleTemplate(
		"proxy/caddy/Caddyfile.tmpl",
		nil,
		nil)
	if err != nil {
		return err
	}
	err = xfs.WriteFile("ssh:///etc/caddy/Caddyfile", caddyFile)
	if err != nil {
		return err
	}

	// Restart caddy
	if _, _, err := system.RemoteRun("systemctl", "restart", "caddy"); err != nil {
		return err
	}

	// Ensure stackhead user can reload caddy
	permissions := "\n%stackhead ALL= NOPASSWD: /bin/systemctl reload caddy\n"
	if err := xfs.AppendToFile("ssh:///etc/sudoers.d/stackhead", permissions, true); err != nil {
		logger.Debugln(err)
		return fmt.Errorf("unable to add Caddy reload permissions for stackhead user")
	}
	// Validate sudoers file
	if _, _, err := system.RemoteRun("/usr/sbin/visudo -cf /etc/sudoers"); err != nil {
		return fmt.Errorf("unable to validate sudoers file")
	}

	// Add stackhead user to docker
	//if _, _, err := system.RemoteRun("usermod", "-a -G docker stackhead"); err != nil {
	//	return fmt.Errorf("unable to add stackhead user to docker group")
	//}

	// adjust owner of /var/www directories
	//if _, _, err := system.RemoteRun("chown", "-R", "stackhead:stackhead", "/var/www"); err != nil {
	//	return err
	//}
	//// adjust owner of /etc/caddy/sites-enabled directories
	//if _, _, err := system.RemoteRun("chown", "-R", "stackhead:stackhead", "/etc/caddy/sites-enabled"); err != nil {
	//	return err
	//}
	//// adjust owner of /etc/caddy/sites-available directories
	//if _, _, err := system.RemoteRun("chown", "-R", "stackhead:stackhead", "/etc/caddy/sites-available"); err != nil {
	//	return err
	//}

	// Check content after provisioning
	//- name: Check content after provisioning
	//  uri:
	//    url: "http://{{ ansible_default_ipv4.address|default(ansible_all_ipv4_addresses[0]) }}"
	//    return_content: yes
	//  register: uri_result
	//  until: '"Caddy web server" in uri_result.content'
	//  retries: 5
	//  delay: 1

	return nil
}

func (Module) Install(modulesSettings interface{}) error {
	return InstallApt()
}
