package routines

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/viper"
)

// Exec is a wrapper function around exec.Command with additional settings for this CLI
func Exec(name string, arg []string) error {
	cmd := exec.Command(name, arg...)
	var outBuffer = new(bytes.Buffer)
	var errBuffer = new(bytes.Buffer)
	if viper.GetBool("verbose") {
		_, err := fmt.Fprintf(os.Stdout, "Executing command: %s\n", strings.Join(append([]string{name}, arg...), " "))
		if err != nil {
			return err
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdout
	} else {
		cmd.Stdout = outBuffer
		cmd.Stderr = errBuffer
	}
	if err := cmd.Run(); err != nil {
		lines := strings.Split(outBuffer.String(), "\n")
		filtered := []string{errBuffer.String()}
		for _, x := range lines {
			if strings.HasPrefix(x, "fatal:") {
				filtered = append(filtered, x)
			}
		}

		return fmt.Errorf(strings.Join(filtered, "\n"))
	}
	return nil
}
