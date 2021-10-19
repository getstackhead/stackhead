package stackhead

import (
	"bytes"
	"fmt"
	"os/exec"
)

func RemoteRun(cmd string) error {
	user := "stackhead"
	if Context.CurrentAction == ContextActionServerSetup {
		user = "root"
	}
	command := exec.Command("ssh", fmt.Sprintf("%s@%s", user, Context.TargetHost), cmd)
	var out bytes.Buffer
	command.Stdout = &out
	return command.Run()
}
