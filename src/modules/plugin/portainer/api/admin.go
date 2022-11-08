package plugin_portainer_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SetupData struct {
	Username string
	Password string
}

func (c Client) InitialSetup(setupData SetupData) (PortainerApiUsersResult, error) {
	var portainerUser PortainerApiUsersResult
	data := map[string]string{
		"Username": setupData.Username,
		"Password": setupData.Password,
	}
	body, _ := json.Marshal(data)
	// post request without auth for admin init
	res, err := c.sendRequest(http.MethodPost, "users/admin/init", portainerUser, bytes.NewBuffer(body), false)
	fmt.Println(res)
	return portainerUser, err
}

func (c Client) GetAuthToken() string {
	data := map[string]string{
		"password": c.AuthPass,
		"username": c.AuthUser,
	}
	body, _ := json.Marshal(data)
	// post request without auth for admin init
	res, err := c.sendRequest(http.MethodPost, "auth", nil, bytes.NewBuffer(body), false)
	if err != nil {
		return ""
	}
	return res["jwt"].(string)
}

// users/admin/check
