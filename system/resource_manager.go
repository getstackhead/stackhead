package system

import (
	"fmt"
	log "github.com/sirupsen/logrus"

	xfs "github.com/saitho/golang-extended-fs/v2"
)

func ApplyResourceOperation(resource *Resource, ignoreBackup bool) (bool, error) {
	if !ignoreBackup {
		// Backup existing file
		backupPath, err := backupResource(resource)
		if err != nil {
			return true, err
		}
		fmt.Println(backupPath)
	}
	return PerformOperation(resource)
}

func RollbackResourceOperation(resource *Resource, ignoreBackup bool) (bool, error) {
	if resource.Operation == OperationCreate {
		resource.Operation = OperationDelete
		found, err := PerformOperation(resource)
		if err != nil {
			return found, err
		}
		if !ignoreBackup {
			// Restore backup
			if err = restoreBackup(resource); err != nil {
				return found, err
			}
		}
		return found, err
	}
	return true, fmt.Errorf(fmt.Sprintf("unupported rollback for operation %s", resource.Operation))
}

func backupResource(resource *Resource) (string, error) {
	//  && resource.Type != TypeLink todo: make it available for symlinks again
	// issue with symlinks: cannot stat symlink: permission denied
	if resource.Type != TypeFile && resource.Type != TypeFolder {
		return "", nil
	}
	if !resource.ExternalResource {
		return "", nil
	}
	resourceFilePath, err := Context.CurrentDeployment.GetResourcePath(resource)
	if err != nil {
		return "", err
	}
	log.Info("Creating backup of resource " + resourceFilePath)
	backupFilePath := resourceFilePath + ".bak"
	xfsFilePath := "ssh://" + resourceFilePath
	switch resource.Type {
	case TypeFile, TypeLink:
		var fileFound bool
		var err error
		if resource.Type == TypeFile {
			fileFound, err = xfs.HasFile(xfsFilePath)
		} else {
			fileFound, err = xfs.HasLink(xfsFilePath)
		}
		if err != nil {
			return "", fmt.Errorf("unable to check status of %s %s: %s", resource.Type, resourceFilePath, err)
		}
		if !fileFound {
			return "", nil
		}
		if _, err = SimpleRemoteRun("cp", RemoteRunOpts{Args: []string{resourceFilePath, backupFilePath}}); err != nil {
			return backupFilePath, fmt.Errorf("unable to backup %s %s: %s", resource.Type, resourceFilePath, err)
		}
		return backupFilePath, nil
	case TypeFolder:
		hasFolder, err := xfs.HasFolder(xfsFilePath)
		if err != nil {
			return "", fmt.Errorf("unable to check status of folder %s: %s", resourceFilePath, err)
		}
		if !hasFolder {
			return "", nil
		}
		if _, err = SimpleRemoteRun("cp", RemoteRunOpts{Args: []string{"-R", resourceFilePath, backupFilePath}}); err != nil {
			return backupFilePath, fmt.Errorf("unable to backup folder %s: %s", resourceFilePath, err)
		}
		return backupFilePath, nil
	}
	return "", fmt.Errorf("unknown backup handler for resource type %s", resource.Type)
}

func restoreBackup(resource *Resource) error {
	if resource.Type != TypeFile && resource.Type != TypeFolder && resource.Type != TypeLink {
		return nil
	}
	if resource.BackupFilePath == "" {
		return nil
	}
	resourceFilePath, _ := Context.CurrentDeployment.GetResourcePath(resource)
	xfsBackupFilePath := "ssh://" + resource.BackupFilePath
	log.Info("Restoring backup of resource " + resourceFilePath)

	switch resource.Type {
	case TypeFile, TypeLink:
		hasFile, err := xfs.HasFile(xfsBackupFilePath)
		if err != nil {
			return err
		}
		if !hasFile {
			return fmt.Errorf("backup not found for " + resource.Name)
		}
		return xfs.CopyFile(xfsBackupFilePath, "ssh://"+resourceFilePath)
	case TypeFolder:
		backupFileName := "ssh://" + resourceFilePath + ".bak"
		hasFolder, err := xfs.HasFolder(backupFileName)
		if err != nil {
			return err
		}
		if !hasFolder {
			return fmt.Errorf("backup not found for " + resource.Name)
		}
		if _, err = SimpleRemoteRun("cp", RemoteRunOpts{Args: []string{"-R", backupFileName, resourceFilePath}}); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("unknown restore backup handler for resource type %s", resource.Type)
}

func PerformOperation(resource *Resource) (bool, error) {
	resourceFilePath, _ := Context.CurrentDeployment.GetResourcePath(resource)
	xfsResourceFilePath := "ssh://" + resourceFilePath
	switch resource.Type {
	case TypeFile:
		if resource.Operation == OperationCreate {
			if err := xfs.WriteFile(xfsResourceFilePath, resource.Content); err != nil {
				return true, fmt.Errorf("unable to create file at %s: %s", resource.Name, err)
			}
		} else if resource.Operation == OperationDelete {
			if err := xfs.DeleteFile(xfsResourceFilePath); err != nil {
				if err.Error() == "file does not exist" {
					return true, nil
				}
				return true, fmt.Errorf("unable to remove file at %s: %s", resource.Name, err)
			}
		}
		return true, nil
	case TypeFolder:
		if resource.Operation == OperationCreate {
			if err := xfs.CreateFolder(xfsResourceFilePath); err != nil {
				return true, fmt.Errorf("unable to create folder at %s: %s", resource.Name, err)
			}
		} else if resource.Operation == OperationDelete {
			if err := xfs.DeleteFolder(xfsResourceFilePath, true); err != nil {
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
			if err := xfs.DeleteFile(xfsResourceFilePath); err != nil {
				if err.Error() == "file does not exist" {
					return true, nil
				}
				return true, fmt.Errorf("unable to remove symlink at %s: %s", resource.Name, err)
			}
		}
		return true, nil
	}
	// CONTAINER via ResourceGroup (see StackHead container module)
	return false, nil
}
