package proxy_nginx

type NginxConfig struct {
	User                      string
	ConfPath                  string
	VhostPath                 string
	ErrorLog                  string
	AccessLog                 string
	PidFile                   string
	WorkerProcesses           string
	ExtraConfOptions          map[string]string
	ExtraConfHttpOptions      map[string]string
	WorkerConnections         int
	MultiAccept               string
	MimeFilePath              string
	ServerNamesHashBucketSize int
	ClientMaxBodySize         string
	Sendfile                  string
	TcpNopush                 string
	TcpNodelay                string
	ServerTokens              string
	ProxyCachePath            string
	KeepaliveTimeout          int
	KeepaliveRequests         int
	TypesHashMaxSize          int
}

func (c *NginxConfig) SetDefaults() {
	if c.User == "" {
		c.User = "www-data"
	}
	if c.ConfPath == "" {
		c.ConfPath = "/etc/nginx/conf.d"
	}
	if c.VhostPath == "" {
		c.VhostPath = "/etc/nginx/sites-enabled"
	}
	if c.ErrorLog == "" {
		c.ErrorLog = "/var/log/nginx/error.log"
	}
	if c.AccessLog == "" {
		c.AccessLog = "/var/log/nginx/access.log"
	}
	if c.PidFile == "" {
		c.PidFile = "/run/nginx.pid"
	}
	if c.WorkerProcesses == "" {
		c.WorkerProcesses = "auto"
	}
	if c.ExtraConfOptions == nil {
		c.ExtraConfOptions = nil
	}
	if c.ExtraConfHttpOptions == nil {
		c.ExtraConfHttpOptions = nil
	}
	if c.WorkerConnections == 0 {
		c.WorkerConnections = 1024
	}
	if c.MultiAccept == "" {
		c.MultiAccept = "off"
	}
	if c.MimeFilePath == "" {
		c.MimeFilePath = "/etc/nginx/mime.types"
	}
	if c.ServerNamesHashBucketSize == 0 {
		c.ServerNamesHashBucketSize = 64
	}
	if c.ClientMaxBodySize == "" {
		c.ClientMaxBodySize = "64m"
	}
	if c.Sendfile == "" {
		c.Sendfile = "on"
	}
	if c.TcpNopush == "" {
		c.TcpNopush = "on"
	}
	if c.TcpNodelay == "" {
		c.TcpNodelay = "on"
	}
	if c.KeepaliveTimeout == 0 {
		c.KeepaliveTimeout = 65
	}
	if c.KeepaliveRequests == 0 {
		c.KeepaliveRequests = 100
	}
	if c.TypesHashMaxSize == 0 {
		c.TypesHashMaxSize = 2048
	}
	if c.ServerTokens == "" {
		c.ServerTokens = "on"
	}
	if c.ProxyCachePath == "" {
		c.ProxyCachePath = ""
	}
}

type ModuleSettings struct {
	CertificatesEmail string
	Config            NginxConfig
}
