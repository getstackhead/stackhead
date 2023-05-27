package system

import (
	"fmt"
	xfs "github.com/saitho/golang-extended-fs/v2"
)

func GetOperationLabel(resource Resource) string {
	operation := "UNKNOWN"
	if resource.Operation == OperationCreate {
		operation = "CREATE"
		if resource.Type == TypeFile {
			// check if file exists
			hasFile, _ := xfs.HasFile("ssh://" + resource.Name)
			if hasFile {
				operation = "UPDATE"
			}
		}
	}
	return operation
}

func ApplyResourcesFromContext() (bool, []error) {
	var resourceRollbackOrder []Resource
	var errors []error
	rollback := false
	for _, resource := range Context.Resources {
		if err := applyResourceOperation(resource); err != nil {
			rollback = true
			errors = append(errors, err)
			break
		}
		resourceRollbackOrder = append([]Resource{resource}, resourceRollbackOrder...)
		if err := resource.ApplyResourceFunc(); err != nil {
			rollback = true
			errors = append(errors, fmt.Errorf("Unable to complete resource creation at %s: %s", resource.Name, err))
			break
		}
	}

	if rollback {
		for _, resource := range resourceRollbackOrder {
			if err := resource.RollbackResourceFunc(); err != nil {
				errors = append(errors, fmt.Errorf("Unable to completely rollback resource at %s: %s", resource.Name, err))
			}
			if err := rollbackResourceOperation(resource); err != nil {
				errors = append(errors, fmt.Errorf("Rollback error: %s", err))
			}
		}
	}
	return !rollback, errors
}
func applyResourceOperation(resource Resource) error {
	// FILE
	if resource.Type == TypeFile {
		if resource.Operation == OperationCreate {
			if err := xfs.WriteFile("ssh://"+resource.Name, resource.Content); err != nil {
				return fmt.Errorf("unable to create file at %s: %s", resource.Name, err)
			}
		}
	}
	// CONTAINER
	// todo
	return nil
}

func rollbackResourceOperation(resource Resource) error {
	// FILE
	if resource.Type == TypeFile {
		if resource.Operation == OperationCreate {
			if err := xfs.DeleteFile("ssh://" + resource.Name); err != nil {
				return fmt.Errorf("unable to remove file at %s: %s", resource.Name, err)
			}
		}
	}
	// CONTAINER
	// todo
	return nil
}
