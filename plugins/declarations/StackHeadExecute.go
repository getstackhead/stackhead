package declarations

import (
	"github.com/getstackhead/stackhead/stackhead"
)

var StackHeadExecute = func(command string) error {
	return stackhead.RemoteRun(command)
}
