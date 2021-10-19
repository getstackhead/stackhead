package stackhead

import (
	"net"

	"github.com/getstackhead/stackhead/pluginlib"
)

var ContextActionProjectDeploy = "project.deploy"
var ContextActionProjectDestroy = "project.destroy"
var ContextActionServerSetup = "server.setup"

type ContextStruct struct {
	TargetHost    net.IP
	CurrentAction string
	Project       *pluginlib.Project
}

var Context = ContextStruct{}

func InitializeContext(host string, action string, projectDefinition *pluginlib.Project) {
	Context.TargetHost = net.IP(host)
	Context.CurrentAction = action
	Context.Project = projectDefinition
}
