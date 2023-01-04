package docker_compose

import (
	"gopkg.in/yaml.v3"
)

type Network struct {
}

type ServiceNetwork struct {
	Aliases []string
}

type DriverOpts struct {
	Type   string
	O      string
	Device string
}

type Volume struct {
	DriverOpts DriverOpts `yaml:"driver_opts,omitempty"`
}

type ServiceVolume struct {
	Type     string
	Source   string
	Target   string
	ReadOnly bool `yaml:"read_only"`
}

type Services struct {
	ContainerName string `yaml:"container_name"`
	Image         string
	Restart       string
	Labels        map[string]string
	User          string                    `yaml:"user,omitempty"`
	Networks      map[string]ServiceNetwork `yaml:"networks,flow"`
	Volumes       []ServiceVolume           `yaml:"volumes,omitempty"`
	Ports         []string                  `yaml:"ports,omitempty"`
	Environment   map[string]string         `yaml:"environment,omitempty"`
	DependsOn     []string                  `yaml:"depends_on,omitempty"`
	VolumesFrom   []string                  `yaml:"volumes_from,omitempty"`
}

type DockerCompose struct {
	Version  string              `yaml:"version"`
	Services map[string]Services `yaml:"services"`
	Networks map[string]Network  `yaml:"networks"`
	Volumes  map[string]Volume   `yaml:"volumes,omitempty"`
}

func (d DockerCompose) String() (string, error) {
	bytes, err := yaml.Marshal(d)
	return string(bytes), err
}
func (d DockerCompose) Map() (map[string]interface{}, error) {
	// unmarshal via yaml, as structs.Map() will not use the keys defined in yaml section but keep the camel case keys
	var unmarshalResult map[string]interface{}
	dcString, err := d.String()
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal([]byte(dcString), &unmarshalResult); err != nil {
		return nil, err
	}
	return unmarshalResult, nil
}

type DockerAuth struct {
	Auth string
}

type LocalDockerConfig struct {
	Auths map[string]DockerAuth
}
