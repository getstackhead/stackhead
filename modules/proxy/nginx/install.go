package proxy_nginx

import (
	"fmt"

	"github.com/fatih/structs"
	xfs "github.com/saitho/golang-extended-fs/v2"
	logger "github.com/sirupsen/logrus"

	"github.com/getstackhead/stackhead/system"
)

func (Module) Install(_modulesSettings interface{}) error {
	moduleSettings, err := system.UnpackModuleSettings[ModuleSettings](_modulesSettings)
	if err != nil {
		return fmt.Errorf("unable to load module settings: " + err.Error())
	}
	moduleSettings.Config.SetDefaults()

	// Add stackhead user to www-data group
	if _, _, err := system.RemoteRun("usermod", system.RemoteRunOpts{Args: []string{"-a -G www-data stackhead"}}); err != nil {
		return fmt.Errorf("unable to add stackhead user to www-data group")
	}

	// Ensure stackhead user can reload nginx
	permissions := "\n%stackhead ALL= NOPASSWD: /bin/systemctl reload nginx\n"
	if err := xfs.AppendToFile("ssh:///etc/sudoers.d/stackhead", permissions, true); err != nil {
		logger.Debugln(err)
		return fmt.Errorf("unable to add Nginx reload permissions for stackhead user")
	}
	// Ensure stackhead user can use certbot
	permissionsCertbot := "\n%stackhead ALL= NOPASSWD: /usr/bin/certbot\n"
	if err := xfs.AppendToFile("ssh:///etc/sudoers.d/stackhead", permissionsCertbot, true); err != nil {
		logger.Debugln(err)
		return fmt.Errorf("unable to add Certbot permissions for stackhead user")
	}
	// Validate sudoers file
	if _, _, err := system.RemoteRun("/usr/sbin/visudo -cf /etc/sudoers", system.RemoteRunOpts{}); err != nil {
		return fmt.Errorf("unable to validate sudoers file")
	}

	if err := system.InstallPackage([]system.Package{
		{
			Name:   "nginx",
			Vendor: system.PackageVendorApt,
		},
		{
			Name:   "certbot",
			Vendor: system.PackageVendorApt,
		},
	}); err != nil {
		return err
	}

	nginxConfig := moduleSettings.Config

	// Override /etc/nginx/nginx.conf
	nginxConfTemplate, err := system.RenderModuleTemplate(
		templates,
		"nginx.conf.tmpl",
		structs.Map(nginxConfig),
		nil)
	if err != nil {
		return err
	}
	err = xfs.WriteFile("ssh:///etc/nginx/nginx.conf", nginxConfTemplate)

	// adjust owner of /var/www directories
	if _, _, err := system.RemoteRun("chown", system.RemoteRunOpts{Args: []string{"-R", "stackhead:stackhead", "/var/www"}}); err != nil {
		return err
	}
	// adjust owner of /etc/nginx/sites-enabled directories
	if _, _, err := system.RemoteRun("chown", system.RemoteRunOpts{Args: []string{"-R", "stackhead:stackhead", moduleSettings.Config.VhostPath}}); err != nil {
		return err
	}
	// adjust owner of /etc/nginx/sites-available directories
	if _, _, err := system.RemoteRun("chown", system.RemoteRunOpts{Args: []string{"-R", "stackhead:stackhead", "/etc/nginx/sites-available"}}); err != nil {
		return err
	}

	// Create certificates folder
	if err := xfs.CreateFolder("ssh://" + CertificatesDirectory); err != nil {
		return err
	}
	if err := xfs.Chown("ssh://"+CertificatesDirectory, 1412, 1412); err != nil {
		return err
	}

	// Create AcmeChallengesDirectory folder
	if err := xfs.CreateFolder("ssh://" + AcmeChallengesDirectory); err != nil {
		return err
	}
	if err := xfs.Chown("ssh://"+AcmeChallengesDirectory, 1412, 1412); err != nil {
		return err
	}

	// Create self-signed certificates
	SnakeoilFullchainPath, SnakeoilPrivkeyPath := GetSnakeoilPaths()
	hasChain, _ := xfs.HasFile("ssh://" + SnakeoilFullchainPath)
	hasKey, _ := xfs.HasFile("ssh://" + SnakeoilPrivkeyPath)
	if !hasChain || !hasKey {
		if _, err := system.SimpleRemoteRun("openssl", system.RemoteRunOpts{
			Args: []string{
				"req",
				"-new",
				"-newkey rsa:4096",
				"-x509",
				"-sha256",
				"-nodes",
				"-out " + SnakeoilFullchainPath,
				"-keyout " + SnakeoilPrivkeyPath,
				"-subj \"/O=StackHead/CN=stackhead.local/\"",
			},
		}); err != nil {
			return fmt.Errorf("Unable to create Snakeoil certifictes: " + err.Error())
		}
	}

	if err := xfs.Chown("ssh://"+SnakeoilFullchainPath, 1412, 1412); err != nil {
		return err
	}
	if err := xfs.Chown("ssh://"+SnakeoilPrivkeyPath, 1412, 1412); err != nil {
		return err
	}

	return nil
}
