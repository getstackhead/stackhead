package plugin_portainer_api

type PortainerApiEndpoint struct {
	Id   string
	Name string
}

func (c Client) GetEndpoints() ([]PortainerApiEndpoint, error) {
	var endpoints []PortainerApiEndpoint
	if _, err := c.get("/endpoints", endpoints); err != nil {
		return endpoints, err
	}
	return endpoints, nil
}
