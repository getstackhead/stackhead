package plugin_portainer_api

func (c Client) CreateTeamMembership(teamId string, userId string) error {
	_, err := c.post(
		"team_memberships",
		map[string]string{
			"role":   "1",
			"teamID": teamId,
			"userID": userId,
		},
		nil,
	)
	return err
}
