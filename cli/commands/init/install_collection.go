package commandsinit

import (
	"fmt"
	"sync"

	"github.com/getstackhead/stackhead/cli/ansible"
	"github.com/getstackhead/stackhead/cli/routines"
)

func installCollection(version string) error {
	repoPath := "git+https://github.com/getstackhead/stackhead.git"
	if len(version) > 0 {
		repoPath += "," + version
	}
	return routines.ExecAnsibleGalaxy(
		"collection", "install", repoPath,
	)
}

func installCollectionDependencies() error {
	collectionDir, err := ansible.GetStackHeadCollectionLocation()
	if err != nil {
		return err
	}
	return routines.ExecAnsibleGalaxy(
		"install",
		"-r", collectionDir+"/requirements/requirements.yml",
	)
}

// InstallCollection is a list of task options that provide the actual workflow being run
func InstallCollection(version string) []routines.TaskOption {
	text := "Installing StackHead Ansible collection"
	if len(version) > 0 {
		text += " (version: " + version + ")"
	}
	return []routines.TaskOption{
		routines.Text(text),
		routines.Execute(func(wg *sync.WaitGroup, result chan routines.TaskResult) {
			defer wg.Done()
			var err error

			// Check if Ansible is installed
			_, err = ansible.GetAnsibleVersion()
			if err != nil {
				err = fmt.Errorf("I could not find Ansible on your system. Please install it")
			}

			if err == nil {
				err = installCollection(version)
			}
			if err == nil {
				err = installCollectionDependencies()
			}

			taskResult := routines.TaskResult{
				Error: err != nil,
			}
			if err != nil {
				taskResult.Message = err.Error()
			}

			result <- taskResult
		}),
	}
}
