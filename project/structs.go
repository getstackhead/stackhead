package project

type Project struct {
	Name                    string
	Domains                 []Domains
	Container               Container
	ProjectDefinitionFolder string
}

type DomainExpose struct {
	InternalPort            int `yaml:"internal_port"`
	ExternalPort            int `yaml:"external_port"`
	Service                 string
	ProxyWebsocketLocations []string `yaml:"proxy_websocket_locations,omitempty"`
}

type DomainSecurityAuthentication struct {
	Type     string
	Username string
	Password string
}

type DomainSecurity struct {
	Authentication []DomainSecurityAuthentication `yaml:"authentication,omitempty"`
}

type Domains struct {
	Domain   string
	Expose   []DomainExpose
	Security DomainSecurity `yaml:"security,omitempty"`
}

type Registries struct {
	Username string
	Password string
	Url      string `yaml:"url,omitempty"`
}

type Container struct {
	Registries []Registries
	Services   []ContainerService
}

type ContainerService struct {
	Name        string
	Image       string
	User        string
	Volumes     []ContainerServiceVolume `yaml:"volumes,omitempty"`
	Hooks       ContainerServiceHooks    `yaml:"hooks,omitempty"`
	VolumesFrom []string                 `yaml:"volumes_from,omitempty"`
	Environment map[string]string        `yaml:"environment,omitempty"`
}

type ContainerServiceHooks struct {
	ExecuteAfterSetup    string `yaml:"execute_after_setup"`
	ExecuteBeforeDestroy string `yaml:"execute_before_destroy"`
}

type ContainerServiceVolume struct {
	Type string
	Src  string
	Dest string
	User string
	Mode string
}
