package proxy_caddy

import (
	"fmt"
	xfs "github.com/saitho/golang-extended-fs/v2"

	"github.com/getstackhead/stackhead/system"
)

func (m Module) Destroy(modulesSettings interface{}) error {
	filePath := "ssh:///etc/caddy/conf.d/stackhead_" + system.Context.Project.Name + ".conf"
	if fileExists, _ := xfs.HasFile(filePath); !fileExists {
		return nil
	}
	if err := xfs.DeleteFile("ssh:///etc/caddy/conf.d/stackhead_" + system.Context.Project.Name + ".conf"); err != nil {
		return fmt.Errorf("Unable to remove symlinked project Caddyfile: " + err.Error())
	}
	if _, err := system.SimpleRemoteRun("systemctl", system.RemoteRunOpts{Args: []string{"reload", "caddy"}, Sudo: true}); err != nil {
		return fmt.Errorf("Unable to reload Caddy service: " + err.Error())
	}
	return nil
}
