package plugin_portainer_api

import "strconv"

type PortainerApiRegistry struct {
	Id   string
	Name string
}

func (c Client) GetRegistries() ([]PortainerApiRegistry, error) {
	var registries []PortainerApiRegistry
	if _, err := c.get("/registries", registries); err != nil {
		return registries, err
	}
	return registries, nil
}

func (c Client) GetRegistry(id int) (*PortainerApiRegistry, error) {
	registries, err := c.GetRegistries()
	if err != nil {
		return nil, err
	}
	for _, registry := range registries {
		if registry.Id == strconv.Itoa(id) {
			return &registry, nil
		}
	}
	return nil, nil
}

func (c Client) CreateRegistry(name string, url string, user string, password string) (PortainerApiRegistry, error) {
	var registry PortainerApiRegistry

	data := map[string]string{
		"authentication": "true",
		"name":           name,
		"url":            url,
		"username":       user,
		"password":       password,
		"type":           "3",
	}

	_, err := c.post("/registries", data, registry)
	if err != nil {
		return registry, err
	}
	return registry, nil
}
