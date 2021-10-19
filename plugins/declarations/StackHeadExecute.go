package declarations

import (
	"bytes"

	"github.com/getstackhead/stackhead/system"
)

var StackHeadExecute = func(command string, args ...string) (bytes.Buffer, bytes.Buffer, error) {
	return system.RemoteRun(command, args...)
}
