package config

import (
	"io/ioutil"
	"os"
	"path"

	yaml "gopkg.in/yaml.v3"

	"github.com/getstackhead/stackhead/pluginlib"
)

var GetPluginDir = func() (string, error) {
	// Find home directory.
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	// Create the file
	return path.Join(configDir, "stackhead", "plugins"), nil
}

func LoadProjectDefinition(filepath string) (*pluginlib.Project, error) {
	p := &pluginlib.Project{}

	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(yamlFile, &p); err != nil {
		return nil, err
	}
	return p, nil
}
