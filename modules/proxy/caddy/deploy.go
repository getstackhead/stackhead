package proxy_caddy

import (
	"fmt"

	xfs "github.com/saitho/golang-extended-fs/v2"

	"github.com/getstackhead/stackhead/modules/proxy"
	"github.com/getstackhead/stackhead/system"
)

func (Module) Deploy(modulesSettings interface{}) error {
	caddyDirectives, err := system.RenderModuleTemplate(
		templates,
		"Caddyfile_project.tmpl",
		nil,
		proxy.FuncMap)
	if err != nil {
		return err
	}

	projectCaddyLocation := system.Context.Project.GetDirectoryPath() + "/Caddyfile"
	if err := xfs.WriteFile("ssh://"+projectCaddyLocation, caddyDirectives); err != nil {
		return err
	}

	if _, err := system.SimpleRemoteRun("ln", system.RemoteRunOpts{Args: []string{"-sf " + projectCaddyLocation + " /etc/caddy/conf.d/stackhead_" + system.Context.Project.Name + ".conf"}}); err != nil {
		return fmt.Errorf("Unable to symlink project Caddyfile: " + err.Error())
	}

	if _, err := system.SimpleRemoteRun("systemctl", system.RemoteRunOpts{Args: []string{"reload", "caddy"}, Sudo: true}); err != nil {
		return fmt.Errorf("Unable to reload Caddy service: " + err.Error())
	}

	return nil
}
