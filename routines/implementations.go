package routines

import (
	"fmt"
	"strings"

	"github.com/chelnak/ysmrr"
	xfs "github.com/saitho/golang-extended-fs/v2"
	logger "github.com/sirupsen/logrus"

	"github.com/getstackhead/stackhead/project"
	"github.com/getstackhead/stackhead/system"
)

var ValidateStackHeadVersionTask = Task{
	Name: fmt.Sprintf("Validating StackHead version"),
	Run: func(r *Task) error {
		isValid, infoText, err := system.ValidateVersion()
		r.SetSuccessMessage(infoText)
		if err != nil {
			logger.Debugln(err)
			err = fmt.Errorf("unable to validate StackHead version.")
		}
		if !isValid {
			err = fmt.Errorf("Trying to deploy with a newer version of StackHead than used for server setup. Please run a server setup again.")
		}

		if err != nil {
			r.SetFailMessage(err.Error())
		}
		return err
	},
}

var PrepareProjectTask = func(projectDefinition *project.Project) Task {
	return Task{
		Name: fmt.Sprintf("Preparing project structure"),
		Run: func(r *Task) error {
			r.PrintLn("Create project directory if not exists")
			if err := xfs.CreateFolder("ssh://" + projectDefinition.GetDirectoryPath()); err != nil {
				return err
			}
			return nil
		},
		ErrorAsErrorMessage: true,
	}
}
var CollectResourcesTask = func(projectDefinition *project.Project) Task {
	return Task{
		Name: fmt.Sprintf("Collecting resources"),
		Run: func(r *Task) error {
			// Collect exposed services
			var exposedServices []project.DomainExpose
			for _, domain := range projectDefinition.Domains {
				exposedServices = append(exposedServices, domain.Expose...)
			}
			for _, module := range system.Context.GetModulesInOrder() {
				if module.GetConfig().Type == "plugin" {
					continue
				}
				moduleSettings := system.GetModuleSettings(module.GetConfig().Name)
				if err := module.Deploy(moduleSettings); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

var resourceRollbackOrder []system.ResourceGroup
var rollback = false

var RollbackResources = Task{
	Name: "Rollback resources",
	Run: func(r *Task) error {
		if !rollback {
			return nil
		}
		var errors []error
		for _, resourceGroup := range resourceRollbackOrder {
			if resourceGroup.RollbackResourceFunc != nil {
				if err := resourceGroup.RollbackResourceFunc(); err != nil {
					errors = append(errors, fmt.Errorf("Unable to completely rollback resources: %s", err))
				}
			}
			for _, resource := range resourceGroup.Resources {
				spinner := r.TaskRunner.GetNewSubtaskSpinner(resource.ToString(true))
				matched, err := system.RollbackResourceOperation(resource)
				if !matched || err == nil {
					spinner.Complete()
				} else if err != nil {
					errors = append(errors, fmt.Errorf("Rollback error: %s", err))
					spinner.Error()
				}
			}
		}
		if len(errors) == 0 {
			return nil
		}
		errorMessages := []string{"The following errors occurred:"}
		for _, err2 := range errors {
			errorMessages = append(errorMessages, "- "+err2.Error())
		}
		return fmt.Errorf(strings.Join(errorMessages, "\n"))
	},
}

var CreateResources = Task{
	Name: "Creating resources",
	Run: func(r *Task) error {
		var errors []error
		var uncompletedSpinners []*ysmrr.Spinner

		for _, resourceGroup := range system.Context.Resources {
			for _, resource := range resourceGroup.Resources {
				spinner := r.TaskRunner.GetNewSubtaskSpinner(resource.ToString(false))
				processed, err := system.ApplyResourceOperation(resource)
				if err != nil {
					rollback = true
					errors = append(errors, err)
					if spinner != nil {
						spinner.UpdateMessage(err.Error())
						spinner.Error()
					}
					return err
				}

				if spinner != nil {
					if processed {
						spinner.Complete()
					} else {
						// uncompleted spinners are resolved when resource group finishes
						uncompletedSpinners = append(uncompletedSpinners, spinner)
					}
				}
			}
			resourceRollbackOrder = append([]system.ResourceGroup{resourceGroup}, resourceRollbackOrder...)
			if resourceGroup.ApplyResourceFunc != nil {
				if err := resourceGroup.ApplyResourceFunc(); err != nil {
					for _, spinner := range uncompletedSpinners {
						spinner.Error()
					}
					rollback = true
					errors = append(errors, fmt.Errorf("Unable to complete resource creation: %s", err))
				}
			}
			if !rollback {
				for _, spinner := range uncompletedSpinners {
					spinner.Complete()
				}
			}
		}
		if !rollback {
			RollbackResources.Disabled = true
			return nil
		}
		if len(errors) == 0 {
			return nil
		}
		errorMessages := []string{"The following errors occurred:"}
		for _, err2 := range errors {
			errorMessages = append(errorMessages, "- "+err2.Error())
		}
		return fmt.Errorf(strings.Join(errorMessages, "\n"))
	},
}
