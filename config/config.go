package config

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

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

	// Set project name. Right now we do not want to allow a "name" attribute in project definition file
	p.Name = strings.TrimRight(path.Base(filepath), ".stackhead.yml")

	return p, nil
}
