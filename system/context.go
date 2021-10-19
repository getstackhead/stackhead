package system

import (
	"net"
	"os"
	"path"

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

func InitializeContext(host string, action string, projectDefinition *pluginlib.Project) {
	Context.TargetHost = net.ParseIP(host)
	Context.CurrentAction = action
	Context.Project = projectDefinition
	Context.IsCI = os.Getenv("CI") != ""
	Context.Authentication = ContextAuthenticationStruct{
		LocalAuthenticationDir: config.GetLocalRemoteKeyDir(Context.TargetHost),
	}
}
