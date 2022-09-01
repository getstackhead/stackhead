package plugin_portainer_api

type PortainerApiUsersResult struct {
	Id       string
	Username string
	Role     int
}

func (c Client) GetUser(username string) (*PortainerApiUsersResult, error) {
	var portainerUsers []PortainerApiUsersResult
	_, err := c.post(
		"users",
		nil,
		portainerUsers,
	)
	if err != nil {
		return nil, err
	}
	for _, user := range portainerUsers {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, nil
}

func (c Client) CreateUser(username string, password string) (PortainerApiUsersResult, error) {
	var portainerUser PortainerApiUsersResult
	_, err := c.post(
		"users",
		map[string]string{
			"username": username,
			"password": password,
			"role":     "2",
		},
		portainerUser,
	)
	if err != nil {
		return portainerUser, err
	}
	return portainerUser, nil
}
