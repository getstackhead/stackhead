package plugin_portainer_api

type PortainerApiTeamResult struct {
	Id   string
	Name string
}

func (c Client) CreateTeam(name string) (PortainerApiTeamResult, error) {
	var portainerTeam PortainerApiTeamResult
	_, err := c.post(
		"teams",
		map[string]string{"name": name},
		portainerTeam,
	)
	if err != nil {
		return portainerTeam, err
	}

	return portainerTeam, nil
}
