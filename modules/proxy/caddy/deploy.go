package proxy_caddy

import (
	"fmt"

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

	caddyFileResource := system.Resource{
		Type:      system.TypeFile,
		Operation: system.OperationCreate,
		Name:      "Caddyfile",
		Content:   caddyDirectives,
	}

	caddyFilePath, err := system.Context.CurrentDeployment.GetResourcePath(&caddyFileResource)
	if err != nil {
		return err
	}

	system.Context.CurrentDeployment.ResourceGroups = append(system.Context.CurrentDeployment.ResourceGroups, system.ResourceGroup{
		Name:      "proxy-caddy-" + system.Context.Project.Name + "-caddyfile",
		Resources: []system.Resource{caddyFileResource},
		ApplyResourceFunc: func() error {
			if _, err := system.SimpleRemoteRun("ln", system.RemoteRunOpts{Args: []string{"-sf " + caddyFilePath + " /etc/caddy/conf.d/stackhead_" + system.Context.Project.Name + ".conf"}}); err != nil {
				return fmt.Errorf("Unable to symlink project Caddyfile: " + err.Error())
			}
			if _, err := system.SimpleRemoteRun("systemctl", system.RemoteRunOpts{Args: []string{"reload", "caddy"}, Sudo: true}); err != nil {
				return fmt.Errorf("Unable to reload Caddy service: " + err.Error())
			}
			return nil
		},
		RollbackResourceFunc: func() error {
			if _, err := system.SimpleRemoteRun("rm", system.RemoteRunOpts{Args: []string{"/etc/caddy/conf.d/stackhead_" + system.Context.Project.Name + ".conf"}}); err != nil {
				return err
			}
			return nil
		},
	})

	return nil
}
