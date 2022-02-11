package system

import (
	logger "github.com/sirupsen/logrus"
	"net"
	"os"
	"path"

	"github.com/saitho/golang-extended-fs/sftp"

	"github.com/getstackhead/stackhead/config"
	"github.com/getstackhead/stackhead/pluginlib"
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
	Project        *pluginlib.Project
	IsCI           bool
	Authentication ContextAuthenticationStruct
}

var Context = ContextStruct{}

type XfsLogger struct {
}

func (l XfsLogger) Debug(obj interface{}) {
	logger.Debug(obj)
}

func (l XfsLogger) Error(obj interface{}) {
	logger.Error(obj)
}

// InitializeContext will set the context object for the current host, action and project
//    host = IP address string
//    action = any of ContextAction* constants
//    projectDefinition = project definition object
func InitializeContext(host string, action string, projectDefinition *pluginlib.Project) {
	Context.TargetHost = net.ParseIP(host)
	Context.CurrentAction = action
	Context.Project = projectDefinition
	Context.IsCI = os.Getenv("CI") != ""
	Context.Authentication = ContextAuthenticationStruct{
		LocalAuthenticationDir: config.GetLocalRemoteKeyDir(Context.TargetHost),
	}

	sftp.Config.Logger = XfsLogger{}
	sftp.Config.SshHost = host
	if action != ContextActionServerSetup {
		// only use private key specifically created during server setup
		sftp.Config.SshUsername = "stackhead"
		sftp.Config.LoadLocalSigners = false
		sftp.Config.SshIdentity = Context.Authentication.GetPrivateKeyPath()
	}
}
