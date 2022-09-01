package plugin_portainer

type PortainerConfig struct {
	ContainerName string `json:"container_name"`
	Server        string
	Port          int
	ApiUser       string `json:"api_user"`
	ApiPassword   string `json:"api_password"`
}

func (p *PortainerConfig) SetDefaults() {
	if p.Server == "" {
		p.Server = "localhost"
	}
	if p.Port == 0 {
		p.Port = 9443
	}
	if p.ContainerName == "" {
		p.ContainerName = "portainer"
	}
}
