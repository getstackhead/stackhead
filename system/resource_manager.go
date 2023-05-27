package system

import (
	"fmt"
	xfs "github.com/saitho/golang-extended-fs/v2"
)

func ApplyResourceOperation(resource Resource) (bool, error) {
	resourceFilePath, _ := Context.CurrentDeployment.GetResourcePath(resource)

	switch resource.Type {
	case TypeFile:
		if resource.Operation == OperationCreate {
			// TODO: backup if file exists
			if err := xfs.WriteFile("ssh://"+resourceFilePath, resource.Content); err != nil {
				return true, fmt.Errorf("unable to create file at %s: %s", resource.Name, err)
			}
		}
		return true, nil
	case TypeFolder:
		if resource.Operation == OperationCreate {
			// TODO: backup if file exists
			if err := xfs.CreateFolder("ssh://" + resourceFilePath); err != nil {
				return true, fmt.Errorf("unable to create folder at %s: %s", resource.Name, err)
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
		}
	}
	// CONTAINER via ResourceGroup (see StackHead container module)
	return false, nil
}

func RollbackResourceOperation(resource Resource) (bool, error) {
	switch resource.Type {
	case TypeFile:
	case TypeLink:
		if resource.Operation == OperationCreate {
			// TODO: restore backup if file exists
			resourcePath, _ := Context.CurrentDeployment.GetResourcePath(resource)
			if err := xfs.DeleteFile("ssh://" + resourcePath); err != nil {
				return true, fmt.Errorf("unable to remove file at %s: %s", resource.Name, err)
			}
		}
		return true, nil
	case TypeFolder:
		if resource.Operation == OperationCreate {
			// TODO: restore backup if file exists
			resourcePath, _ := Context.CurrentDeployment.GetResourcePath(resource)
			if err := xfs.DeleteFolder("ssh://"+resourcePath, true); err != nil {
				return true, fmt.Errorf("unable to remove folder at %s: %s", resource.Name, err)
			}
		}
		return true, nil
	}
	// CONTAINER via ResourceGroup (see StackHead container module)
	return false, nil
}
