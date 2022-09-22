package system

import (
	logger "github.com/sirupsen/logrus"
	"net"
	"os"
	"path"

	"github.com/saitho/golang-extended-fs/v2/sftp"

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
	DNSModules      []Module
	PluginModules   []Module
}

func (c ContextStruct) GetModulesInOrder() []Module {
	modules := []Module{c.ContainerModule, c.ProxyModule}
	modules = append(modules, c.DNSModules...)
	modules = append(modules, c.PluginModules...)
	return modules
}

var Context = ContextStruct{}

type DebugLogger struct {
}

func (DebugLogger) Debug(obj interface{}) {
	logger.Debugln(obj)
}
func (DebugLogger) Error(obj interface{}) {
	logger.Errorln(obj)
}

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
	sftp.Config.Logger = DebugLogger{}
}

func ContextSetProxyModule(module Module) {
	if module.GetConfig().Type != "proxy" {
		return
	}
	Context.ProxyModule = module
}

func ContextSetContainerModule(module Module) {
	if module.GetConfig().Type != "container" {
		return
	}
	Context.ContainerModule = module
}

func ContextAddDnsModule(module Module) {
	if module.GetConfig().Type != "dns" {
		return
	}
	Context.DNSModules = append(Context.DNSModules, module)
}

func ContextAddPluginModule(module Module) {
	if module.GetConfig().Type != "plugin" {
		return
	}
	Context.PluginModules = append(Context.PluginModules, module)
}
