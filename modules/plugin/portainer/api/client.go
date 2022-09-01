package plugin_portainer_api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Client struct {
	Host string
	Port int

	AuthUser string
	AuthPass string
}

func (c Client) getHttpClient() *http.Client {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return &http.Client{Transport: customTransport}
}

func parseResponse(resp *http.Response, resObj *interface{}) (map[string]interface{}, error) {
	defer resp.Body.Close()
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resObj != nil {
		if err := json.Unmarshal(res, &resObj); err != nil {
			return nil, err
		}
	}
	result := make(map[string]interface{})
	err = json.Unmarshal(res, &result)

	if details, ok := result["details"]; ok {
		if details == "Unauthorized" {
			if message, ok := result["message"]; ok {
				return nil, fmt.Errorf(message.(string))
			} else {
				return nil, fmt.Errorf("An API error occurred (" + details.(string) + ")")
			}
		}
	}

	return result, err
}

func (c Client) sendRequest(method string, path string, resObj interface{}, body io.Reader, requireAuth bool) (map[string]interface{}, error) {
	req, err := http.NewRequest(method, c.getApiUrl()+path, body)
	if err != nil {
		return nil, err
	}
	resp, err := c.getHttpClient().Do(req)
	req.Header.Set("Content-Type", "application/json")
	if requireAuth {
		req.Header.Set("Authorization", "Bearer "+c.GetAuthToken())
	}
	if err != nil {
		return nil, err
	}
	return parseResponse(resp, &resObj)
}

func (c Client) get(path string, resObj interface{}) (map[string]interface{}, error) {
	return c.sendRequest(http.MethodGet, path, resObj, nil, true)
}

func (c Client) post(path string, data map[string]string, resObj interface{}) (map[string]interface{}, error) {
	body, _ := json.Marshal(data)
	return c.sendRequest(http.MethodPost, path, resObj, bytes.NewBuffer(body), true)
}

func (c Client) put(path string, data map[string]string, resObj interface{}) (map[string]interface{}, error) {
	// TODO: Auth
	body, _ := json.Marshal(data)

	return c.sendRequest(http.MethodPut, path, resObj, bytes.NewBuffer(body), true)
}

func (c Client) delete(path string) (map[string]interface{}, error) {
	return c.sendRequest(http.MethodDelete, path, nil, bytes.NewBuffer(nil), true)
}

func (c Client) getApiUrl() string {
	return "https://" + c.Host + ":" + strconv.Itoa(c.Port) + "/api/"
}
