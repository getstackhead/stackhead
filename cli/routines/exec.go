package routines

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	"github.com/getstackhead/stackhead/cli/ansible"
)

// Exec is a wrapper function around exec.Command with additional settings for this CLI
func Exec(name string, arg ...string) error {
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

// ExecAnsibleGalaxy is shortcut for executing a command via ansible-galaxy binary
func ExecAnsibleGalaxy(args ...string) error {
	collectionDir, err := ansible.GetCollectionDirs()
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// We have to set a relative path, otherwise ansible-galaxy will do weird things...
	relCollectionsPath, _ := filepath.Rel(cwd, "/")
	args = append(args, "-p "+relCollectionsPath+"/../.."+collectionDir[0])
	return Exec("ansible-galaxy", args...)
}

// ExecAnsiblePlaybook is shortcut for executing a playbook within the StackHead collection via ansible-playbook binary
func ExecAnsiblePlaybook(playbookName string, inventoryPath string, options map[string]string) error {
	stackHeadLocation, err := ansible.GetStackHeadCollectionLocation()
	if err != nil {
		return err
	}

	args := []string{stackHeadLocation + "/playbooks/" + playbookName + ".yml"}
	if inventoryPath != "" {
		args = append(args, "-i", inventoryPath)
	}
	if len(options) > 0 {
		var extraVars []string
		for i, arg := range options {
			extraVars = append(extraVars, i+"="+arg)
		}
		args = append(args, "--extra-vars", strings.Join(extraVars, ","))
	}

	return Exec("ansible-playbook", args...)
}
