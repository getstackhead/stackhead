package project

import (
	"path"
	"strings"

	xfs "github.com/saitho/golang-extended-fs/v2"
	"gopkg.in/yaml.v3"
)

func LoadProjectDefinition(filepath string) (*Project, error) {
	p := &Project{}

	yamlFile, err := xfs.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal([]byte(yamlFile), &p); err != nil {
		return nil, err
	}

	// Set project name. Right now we do not want to allow a "name" attribute in project definition file
	p.Name = strings.TrimRight(path.Base(filepath), ".stackhead.yml")

	return p, nil
}
