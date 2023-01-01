package plugin_portainer

import (
	"fmt"
	"github.com/gookit/event"
	logger "github.com/sirupsen/logrus"

	plugin_portainer_api "github.com/getstackhead/stackhead/modules/plugin/portainer/api"
	"github.com/getstackhead/stackhead/system"
)

func (Module) Init(moduleSettings interface{}) {
	event.On("setup.modules.post-install-module.container.docker", event.ListenerFunc(func(e event.Event) error {
		if _, _, err := system.RemoteRun("docker", system.RemoteRunOpts{Args: []string{"volume", "create", "portainer_data"}}); err != nil {
			return err
		}

		var err error
		var portainerConfig *PortainerConfig
		if moduleSettings == nil {
			portainerConfig = &PortainerConfig{}
		} else {
			portainerConfig, err = system.UnpackModuleSettings[PortainerConfig](moduleSettings)
			if err != nil {
				return err
			}
		}
		portainerConfig.SetDefaults()

		if portainerConfig.ApiUser == "" {
			return fmt.Errorf("missing Portainer API user. please set one on CLI config")
		}
		if portainerConfig.ApiPassword == "" {
			return fmt.Errorf("missing Portainer API password. please set one on CLI config")
		}

		host := portainerConfig.Server
		if host == "localhost" {
			host = system.Context.TargetHost.String()
		}

		portainerApi := plugin_portainer_api.Client{
			Host:     host,
			Port:     portainerConfig.Port,
			AuthUser: portainerConfig.ApiUser,
			AuthPass: portainerConfig.ApiPassword,
		}

		// Setup Portainer container
		//var portainerUser plugin_portainer_api.PortainerApiUsersResult
		result, _, err := system.RemoteRun("docker", system.RemoteRunOpts{Args: []string{"ps", "-a --format {{.Names}} | grep " + portainerConfig.ContainerName + " -w"}})
		if err != nil || result.Len() == 0 {
			if _, errMsg, err := system.RemoteRun(
				"docker",
				system.RemoteRunOpts{
					Args: []string{
						"run", "-d",
						"-p", "8000:8000", "-p", "9443:9443",
						"--name", portainerConfig.ContainerName,
						"--restart=always",
						"-v", "/var/run/docker.sock:/var/run/docker.sock",
						"-v", "portainer_data:/data",
						"portainer/portainer-ce:2.11.1",
					},
				},
			); err != nil {
				logger.Errorln(errMsg.String())
				return err
			}
			// Initial Portainer setup via HTTP request
			setupData := plugin_portainer_api.SetupData{
				Username: portainerConfig.ApiUser,
				Password: portainerConfig.ApiPassword,
			}
			_, err = portainerApi.InitialSetup(setupData)
			if err != nil {
				logger.Debugln("Error setting up initial portainer admin.")
				return err
			}
		} else {
			logger.Debugln("Portainer container \"" + portainerConfig.ContainerName + "\" is already running.")
			// Find portainerUser
			//_, err := portainerApi.GetUser(portainerConfig.ApiUser)
			//if err != nil {
			//	return err
			//}
		}

		// todo: check if team is needed

		//portainerTeam, err := portainerApi.CreateTeam("stackhead")
		//if err != nil {
		//	return err
		//}
		//
		//if err := portainerApi.CreateTeamMembership(portainerTeam.Id, portainerUser.Id); err != nil {
		//	return err
		//}
		return nil
	}), event.Normal)
}
