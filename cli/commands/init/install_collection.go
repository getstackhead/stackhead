package commands_init

import (
	"sync"

	"github.com/getstackhead/stackhead/cli/ansible"
	"github.com/getstackhead/stackhead/cli/routines"
)

func installCollection() error {
	return routines.ExecAnsibleGalaxy(
		"collection", "install", "git+https://github.com/getstackhead/stackhead.git",
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

var InstallCollection = []routines.TaskOption{
	routines.Text("Installing StackHead Ansible collection"),
	routines.Execute(func(wg *sync.WaitGroup, result chan routines.TaskResult) {
		defer wg.Done()

		err := installCollection()
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
