package declarations

import (
	"github.com/getstackhead/stackhead/pluginlib"
	"github.com/getstackhead/stackhead/system"
)

var GetProject = func() *pluginlib.Project {
	return system.Context.Project
}
