package proxy_nginx

import (
	"fmt"
	xfs "github.com/saitho/golang-extended-fs/v2"
	"path"
	"strings"

	"github.com/getstackhead/stackhead/system"
)

func (m Module) Destroy(_modulesSettings interface{}) error {
	moduleSettings, err := system.UnpackModuleSettings[ModuleSettings](_modulesSettings)
	if err != nil {
		return fmt.Errorf("unable to load module settings: " + err.Error())
	}
	moduleSettings.Config.SetDefaults()

	firstDomain := system.Context.Project.Domains[0].Domain
	if _, err := system.SimpleRemoteRun("certbot", system.RemoteRunOpts{Args: []string{"delete", "-q", "--cert-name " + firstDomain}, Sudo: true}); err != nil {
		if !strings.Contains(err.Error(), "No certificate found") {
			return fmt.Errorf("Unable to remove Certbot certificate: " + err.Error())
		}
	}

	domainChallengeDir := path.Join(AcmeChallengesDirectory, firstDomain)
	if err := xfs.DeleteFolder("ssh://"+domainChallengeDir, true); err != nil && err.Error() != "file does not exist" {
		return fmt.Errorf("Unable to remove ACME challenge directory: " + err.Error())
	}

	if _, err := system.SimpleRemoteRun("systemctl", system.RemoteRunOpts{Args: []string{"reload", "nginx"}, Sudo: true}); err != nil {
		return fmt.Errorf("Unable to reload Nginx service: " + err.Error())
	}

	return nil
}
