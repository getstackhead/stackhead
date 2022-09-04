package proxy_nginx

import (
	"fmt"

	"github.com/fatih/structs"
	xfs "github.com/saitho/golang-extended-fs"
	logger "github.com/sirupsen/logrus"

	"github.com/getstackhead/stackhead/system"
)

func (Module) Install(_modulesSettings interface{}) error {
	moduleSettings, err := system.UnpackModuleSettings[ModuleSettings](_modulesSettings)
	if err != nil {
		return fmt.Errorf("unable to load module settings: " + err.Error())
	}
	moduleSettings.Config.SetDefaults()

	// Ensure stackhead user can reload nginx
	permissions := "\n%stackhead ALL= NOPASSWD: /bin/systemctl reload nginx\n"
	if err := xfs.AppendToFile("ssh:///etc/sudoers.d/stackhead", permissions, true); err != nil {
		logger.Debugln(err)
		return fmt.Errorf("unable to add Nginx reload permissions for stackhead user")
	}
	// Validate sudoers file
	if _, _, err := system.RemoteRun("/usr/sbin/visudo -cf /etc/sudoers"); err != nil {
		return fmt.Errorf("unable to validate sudoers file")
	}

	if err := system.InstallPackage([]system.Package{
		{
			Name:   "nginx",
			Vendor: system.PackageVendorApt,
		},
	}); err != nil {
		return err
	}

	nginxConfig := moduleSettings.Config

	// Override /etc/nginx/nginx.conf
	nginxConfTemplate, err := system.RenderModuleTemplate(
		"proxy/nginx/nginx.conf.tmpl",
		structs.Map(nginxConfig),
		nil)
	if err != nil {
		return err
	}
	err = xfs.WriteFile("ssh:///etc/nginx/nginx.conf", nginxConfTemplate)

	// adjust owner of /var/www directories
	if _, _, err := system.RemoteRun("chown", "-R", "stackhead:stackhead", "/var/www"); err != nil {
		return err
	}
	// adjust owner of /etc/nginx/sites-enabled directories
	if _, _, err := system.RemoteRun("chown", "-R", "stackhead:stackhead", "/etc/nginx/sites-enabled"); err != nil {
		return err
	}
	// adjust owner of /etc/nginx/sites-available directories
	if _, _, err := system.RemoteRun("chown", "-R", "stackhead:stackhead", "/etc/nginx/sites-available"); err != nil {
		return err
	}

	// Create certificates folder
	if err := xfs.CreateFolder("ssh://" + CertificatesDirectory); err != nil {
		return err
	}
	if err := xfs.Chown("ssh://"+CertificatesDirectory, 1412, 1412); err != nil {
		return err
	}

	SnakeoilFullchainPath, SnakeoilPrivkeyPath := GetSnakeoilPaths()
	if err := xfs.Chown("ssh://"+SnakeoilFullchainPath, 1412, 1412); err != nil {
		return err
	}
	if err := xfs.Chown("ssh://"+SnakeoilPrivkeyPath, 1412, 1412); err != nil {
		return err
	}

	// Check content after provisioning
	//- name: Check content after provisioning
	//  uri:
	//    url: "http://{{ ansible_default_ipv4.address|default(ansible_all_ipv4_addresses[0]) }}"
	//    return_content: yes
	//  register: uri_result
	//  until: '"Welcome to nginx" in uri_result.content'
	//  retries: 5
	//  delay: 1
	return nil
}
