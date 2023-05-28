package system

import (
	"fmt"

	xfs "github.com/saitho/golang-extended-fs/v2"
)

func ApplyResourceOperation(resource Resource, ignoreBackup bool) (bool, error) {
	return PerformOperation(resource, ignoreBackup)
}

func RollbackResourceOperation(resource Resource, ignoreBackup bool) (bool, error) {
	if resource.Operation == OperationCreate {
		resource.Operation = OperationDelete
		return PerformOperation(resource, ignoreBackup)
	}
	return true, fmt.Errorf(fmt.Sprintf("unupported rollback for operation %s", resource.Operation))
}

func PerformOperation(resource Resource, ignoreBackup bool) (bool, error) {
	resourceFilePath, _ := Context.CurrentDeployment.GetResourcePath(resource)
	switch resource.Type {
	case TypeFile:
		if resource.Operation == OperationCreate {
			// TODO: backup if file exists
			if err := xfs.WriteFile("ssh://"+resourceFilePath, resource.Content); err != nil {
				return true, fmt.Errorf("unable to create file at %s: %s", resource.Name, err)
			}
		} else if resource.Operation == OperationDelete {
			// TODO: restore backup if file exists
			resourcePath, _ := Context.CurrentDeployment.GetResourcePath(resource)
			if err := xfs.DeleteFile("ssh://" + resourcePath); err != nil {
				if err.Error() == "file does not exist" {
					return true, nil
				}
				return true, fmt.Errorf("unable to remove file at %s: %s", resource.Name, err)
			}
		}
		return true, nil
	case TypeFolder:
		if resource.Operation == OperationCreate {
			// TODO: backup if file exists
			if err := xfs.CreateFolder("ssh://" + resourceFilePath); err != nil {
				return true, fmt.Errorf("unable to create folder at %s: %s", resource.Name, err)
			}
		} else if resource.Operation == OperationDelete {
			// TODO: restore backup if file exists
			if err := xfs.DeleteFolder("ssh://"+resourceFilePath, true); err != nil {
				return true, fmt.Errorf("unable to remove folder at %s: %s", resource.Name, err)
			}
		}
		return true, nil
	case TypeLink:
		if resource.Operation == OperationCreate {
			args := RemoteRunOpts{Args: []string{"-s " + resource.LinkSource + " " + resourceFilePath}, AllowFail: true}
			if resource.EnforceLink {
				args = RemoteRunOpts{Args: []string{"-sf " + resource.LinkSource + " " + resourceFilePath}}
			}
			if _, err := SimpleRemoteRun("ln", args); err != nil {
				return true, fmt.Errorf("Unable to symlink " + resource.LinkSource + " -> " + resourceFilePath + ": " + err.Error())
			}
		} else if resource.Operation == OperationDelete {
			// TODO: restore backup if file exists
			if err := xfs.DeleteFile("ssh://" + resourceFilePath); err != nil {
				if err.Error() == "file does not exist" {
					return true, nil
				}
				return true, fmt.Errorf("unable to remove symlink at %s: %s", resource.Name, err)
			}
		}
	}
	// CONTAINER via ResourceGroup (see StackHead container module)
	return false, nil
}
