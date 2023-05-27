package system

import (
	"fmt"

	xfs "github.com/saitho/golang-extended-fs/v2"
)

func ApplyResourceOperation(resource Resource) (bool, error) {
	// FILE
	if resource.Type == TypeFile {
		if resource.Operation == OperationCreate {
			// TODO: backup if file exists
			if err := xfs.WriteFile("ssh://"+resource.Name, resource.Content); err != nil {
				return true, fmt.Errorf("unable to create file at %s: %s", resource.Name, err)
			}
		}
		return true, nil
	}
	// CONTAINER via ResourceGroup (see StackHead container module)
	return false, nil
}

func RollbackResourceOperation(resource Resource) (bool, error) {
	// FILE
	if resource.Type == TypeFile {
		if resource.Operation == OperationCreate {
			// TODO: restore backup if file exists
			if err := xfs.DeleteFile("ssh://" + resource.Name); err != nil {
				return true, fmt.Errorf("unable to remove file at %s: %s", resource.Name, err)
			}
		}
		return true, nil
	}
	// CONTAINER via ResourceGroup (see StackHead container module)
	return false, nil
}
