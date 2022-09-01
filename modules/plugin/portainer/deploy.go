package plugin_portainer

import (
	"fmt"

	plugin_portainer_api "github.com/getstackhead/stackhead/modules/plugin/portainer/api"
	"github.com/getstackhead/stackhead/system"
)

func (m Module) Deploy(moduleSettings interface{}) error {
	// NOT CALLED as this is a plugin
	// todo: use Portainer to add containers when target server is a thing

	portainerConfig, err := system.UnpackModuleSettings[PortainerConfig](moduleSettings)
	if err != nil {
		return err
	}
	portainerConfig.SetDefaults()
	host := portainerConfig.Server
	if host == "localhost" {
		host = system.Context.TargetHost.String()
	}

	portainerApi := plugin_portainer_api.Client{Host: host, Port: portainerConfig.Port}

	//project := system.Context.Project

	// add registry credentials to Portainer

	// /var/lib/docker/volumes/portainer_data/_data/compose

	// GET https://167.235.77.57:9443/api/custom_templates?type=2
	//
	// POST https://167.235.77.57:9443/api/custom_templates
	// method=string
	//  body_string
	//
	// {
	//  "description": "High performance web server",
	//  "fileContent": "string",
	//  "logo": "https://cloudinovasi.id/assets/img/logos/nginx.png",
	//  "note": "This is my <b>custom</b> template",
	//  "platform": 1,
	//  "title": "Nginx",
	//  "type": 2,
	//  "variables": [
	//    {
	//      "defaultValue": "default value",
	//      "description": "Description",
	//      "label": "My Variable",
	//      "name": "MY_VAR"
	//    }
	//  ]
	//}
	//

	endpoints, err := portainerApi.GetEndpoints()
	if err != nil {
		return fmt.Errorf("unable to check Portainer endpoints: " + err.Error())
	}
	if len(endpoints) == 0 {
		return fmt.Errorf("No Endpoints configured on Portainer")
	}
	if len(endpoints) > 1 {
		// todo: allow multiple endpoints; requires selection of one
		return fmt.Errorf("More than one endpoints configured on Portainer")
	}
	currentEndpoint := endpoints[0]
	fmt.Println(currentEndpoint)

	// POST => create
	//err := portainerApi.CreateStack(123, "...body...")

	// PUT => update
	//err := portainerApi.UpdateStack(1, "...body...")

	// DELETE => remove
	//err := portainerApi.DeleteStack(1)

	// POST /stacks
	/**
	{
		"type": 2,
		"method": "string",
	  "endpointId": "123",
		"body_compose_string": "..."
	}
	*/

	//err := portainerApi.CreateRegistry("stackhead-APPNAME", "url", "user", "password")

	return nil
}
