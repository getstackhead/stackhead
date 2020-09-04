package routines

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/viper"

	"github.com/getstackhead/stackhead/cli/ansible"
)

func Exec(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	var errBuffer = new(bytes.Buffer)
	if viper.GetBool("verbose") {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("Executing command: %s", strings.Join(append([]string{name}, arg...), " ")))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdout
	} else {
		cmd.Stderr = errBuffer
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf(errBuffer.String())
	}
	return nil
}

func ExecAnsibleGalaxy(args ...string) error {
	collectionDir, err := ansible.GetCollectionDirs()
	if err != nil {
		return err
	}
	args = append(args, "-p "+collectionDir[0])
	return Exec("ansible-galaxy", args...)
}
