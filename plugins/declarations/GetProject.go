package declarations

import (
	"github.com/getstackhead/stackhead/pluginlib"
	"github.com/getstackhead/stackhead/stackhead"
)

var GetProject = func() *pluginlib.Project {
	return stackhead.Context.Project
}
