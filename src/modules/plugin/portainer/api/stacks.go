package plugin_portainer_api

import (
	"strconv"
)

type PortainerApiStack struct {
	Id   string
	Name string
}

func (c Client) CreateStack(endpointId int, body string) (PortainerApiStack, error) {
	var stack PortainerApiStack

	data := map[string]string{
		"type":                "2",
		"method":              "string",
		"endpointId":          strconv.Itoa(endpointId),
		"body_compose_string": body,
	}

	_, err := c.post("/stacks", data, &stack)
	if err != nil {
		return stack, err
	}
	return stack, nil
}

func (c Client) UpdateStack(stackId int, body string) (PortainerApiStack, error) {
	var stack PortainerApiStack

	data := map[string]string{
		//"endpointId":      strconv.Itoa(endpointId), // todo: does changing the endpoint move the container + data to a different server??
		"body_compose_string": body,
	}

	_, err := c.put("/stacks/"+strconv.Itoa(stackId), data, &stack)
	if err != nil {
		return stack, err
	}
	return stack, nil
}

func (c Client) DeleteStack(stackId int) error {
	if _, err := c.delete("/stacks/" + strconv.Itoa(stackId)); err != nil {
		return err
	}
	return nil
}
