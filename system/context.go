package system

import (
	"net"
	"os"
	"path"

	"github.com/saitho/golang-extended-fs/sftp"

	"github.com/getstackhead/stackhead/config"
	"github.com/getstackhead/stackhead/project"
)

var ContextActionProjectDeploy = "project.deploy"
var ContextActionProjectDestroy = "project.destroy"
var ContextActionServerSetup = "server.setup"

type ContextAuthenticationStruct struct {
	LocalAuthenticationDir string
}

func (c ContextAuthenticationStruct) GetPrivateKeyPath() string {
	return path.Join(c.LocalAuthenticationDir, "private_key.pem")
}

func (c ContextAuthenticationStruct) GetPublicKeyPath() string {
	return path.Join(c.LocalAuthenticationDir, "public_key.pem")
}

type ContextStruct struct {
	TargetHost     net.IP
	CurrentAction  string
	Project        *project.Project
	IsCI           bool
	Authentication ContextAuthenticationStruct

	ProxyModule     Module
	ContainerModule Module
}

func (c ContextStruct) GetModulesInOrder() []Module {
	return []Module{c.ContainerModule, c.ProxyModule}
}

var Context = ContextStruct{}

// InitializeContext will set the context object for the current host, action and project
//    host = IP address string
//    action = any of ContextAction* constants
//    projectDefinition = project definition object
func InitializeContext(host string, action string, projectDefinition *project.Project) {
	Context.TargetHost = net.ParseIP(host)
	Context.CurrentAction = action
	Context.Project = projectDefinition
	Context.IsCI = os.Getenv("CI") != ""
	Context.Authentication = ContextAuthenticationStruct{
		LocalAuthenticationDir: config.GetLocalRemoteKeyDir(Context.TargetHost),
	}

	sftp.Config.SshHost = host
	if action != ContextActionServerSetup {
		// only use private key specifically created during server setup
		sftp.Config.SshUsername = "stackhead"
		sftp.Config.LoadLocalSigners = false
		sftp.Config.SshIdentity = Context.Authentication.GetPrivateKeyPath()
	}
}

func ContextSetProxyModule(proxyModule Module) {
	Context.ProxyModule = proxyModule
}

func ContextSetContainerModule(containerModule Module) {
	Context.ContainerModule = containerModule
}
