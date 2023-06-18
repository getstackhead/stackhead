package routines

import (
	"fmt"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/chelnak/ysmrr"
	xfs "github.com/saitho/golang-extended-fs/v2"
	logger "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

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
		Name: fmt.Sprintf("Preparing deployment"),
		Run: func(r *Task) error {
			r.PrintLn("Create project directory if not exists")
			if err := xfs.CreateFolder("ssh://" + projectDefinition.GetDirectoryPath()); err != nil {
				return err
			}
			if err := xfs.CreateFolder("ssh://" + projectDefinition.GetDeploymentsPath()); err != nil {
				return err
			}

			r.PrintLn("Lookup previous deployments")
			// Find latest deployment
			latestDeployment, err := system.GetLatestDeployment(projectDefinition)
			if err != nil {
				return err
			}
			system.Context.LatestDeployment = latestDeployment
			oldVersion := "N/A"
			newVersion := 1
			if system.Context.LatestDeployment != nil {
				oldVersion = "v" + strconv.Itoa(system.Context.LatestDeployment.Version)
				newVersion = system.Context.LatestDeployment.Version + 1
			}
			system.Context.CurrentDeployment = system.Deployment{
				Version:   newVersion,
				DateStart: time.Now(),
				Project:   system.Context.Project,
			}
			r.PrintLn(fmt.Sprintf("Previous deployment: %s, new deployment: v%d", oldVersion, newVersion))

			// Create folder for new deployment
			if err := xfs.CreateFolder("ssh://" + system.Context.CurrentDeployment.GetPath()); err != nil {
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
				r.PrintLn("Collecting from " + module.GetConfig().Name)
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
			if _, err := processResourceGroup(r.TaskRunner, resourceGroup, true, false); err != nil {
				errors = append(errors, fmt.Errorf("Rollback error: %s", err))
			}
		}

		// Mark deployment as rolled back
		system.Context.CurrentDeployment.RolledBack = true
		for _, err2 := range errors {
			system.Context.CurrentDeployment.RollbackErrors = append(system.Context.CurrentDeployment.RollbackErrors, err2.Error())
		}

		if len(system.Context.CurrentDeployment.RollbackErrors) > 0 {
			return fmt.Errorf("The following errors occurred:\n" + strings.Join(system.Context.CurrentDeployment.RollbackErrors, "\n"))
		}

		return nil
	},
}

// return: bool: whether to consider resource group for requiring rollback ; error
func processResourceGroup(taskRunner *TaskRunner, resourceGroup system.ResourceGroup, isRollbackMode bool, ignoreBackup bool) (bool, error) {
	var uncompletedSpinners []*ysmrr.Spinner

	// ROLLBACK mode
	if isRollbackMode && resourceGroup.RollbackResourceFunc != nil {
		if err := resourceGroup.RollbackResourceFunc(); err != nil {
			return false, err
		}
	}

	for _, resource := range resourceGroup.Resources {
		spinner := taskRunner.GetNewSubtaskSpinner(resource.ToString(isRollbackMode))
		var err error
		var processed bool
		if isRollbackMode {
			processed, err = system.RollbackResourceOperation(&resource, ignoreBackup)
		} else {
			processed, err = system.ApplyResourceOperation(&resource, ignoreBackup)
		}
		if err != nil {
			if spinner != nil {
				spinner.UpdateMessage(err.Error())
				spinner.Error()
			}
			return false, err
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

	// APPLY mode
	if !isRollbackMode && resourceGroup.ApplyResourceFunc != nil {
		if err := resourceGroup.ApplyResourceFunc(); err != nil {
			for _, spinner := range uncompletedSpinners {
				spinner.Error()
			}
			return true, err
		}
	}
	for _, spinner := range uncompletedSpinners {
		spinner.Complete()
	}
	return !isRollbackMode, nil
}

var CreateResources = Task{
	Name: "Creating resources",
	Run: func(r *Task) error {
		var errors []string
		for _, resourceGroup := range system.Context.CurrentDeployment.ResourceGroups {
			considerForRollback, err := processResourceGroup(r.TaskRunner, resourceGroup, false, false)
			if considerForRollback {
				resourceRollbackOrder = append([]system.ResourceGroup{resourceGroup}, resourceRollbackOrder...)
			}
			if err != nil {
				rollback = true
				errors = append(errors, err.Error())
				break
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
			errorMessages = append(errorMessages, "- "+err2)
		}
		return fmt.Errorf(strings.Join(errorMessages, "\n"))
	},
}

func reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

var RemoveResources = func(latestDeployment *system.Deployment) Task {
	return Task{
		Name: "Removing project resources",
		Run: func(r *Task) error {
			reverse(latestDeployment.ResourceGroups)
			for _, group := range latestDeployment.ResourceGroups {
				var filteredResources []system.Resource
				for _, resource := range group.Resources {
					if resource.ExternalResource {
						resource.Operation = system.OperationDelete
						filteredResources = append(filteredResources, resource)
					}
				}
				reverse(filteredResources)
				group.Resources = filteredResources

				if _, err := processResourceGroup(r.TaskRunner, group, false, true); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

var FinalizeDeployment = Task{
	Name: "Finalizing deployment",
	Run: func(r *Task) error {
		// set deployment end date
		system.Context.CurrentDeployment.DateEnd = time.Now()

		// save deployment.yaml file
		yamlString, err := yaml.Marshal(system.Context.CurrentDeployment)
		if err != nil {
			return err
		}
		if err = xfs.WriteFile("ssh://"+path.Join(system.Context.CurrentDeployment.GetPath(), "deployment.yaml"), string(yamlString)); err != nil {
			return err
		}

		if !system.Context.CurrentDeployment.RolledBack {
			// Remove external backups
			for _, resourceGroup := range system.Context.CurrentDeployment.ResourceGroups {
				for _, resource := range resourceGroup.Resources {
					fmt.Println(resource.BackupFilePath) // todo: remove
				}
			}

			// update current symlink if deployment was successful
			if _, err := system.SimpleRemoteRun("ln", system.RemoteRunOpts{Args: []string{"-sfn " + system.Context.CurrentDeployment.GetPath() + " " + path.Join(system.Context.CurrentDeployment.Project.GetDeploymentsPath(), "current")}}); err != nil {
				return fmt.Errorf("Unable to symlink current deployment: " + err.Error())
			}
		}
		return nil
	},
}
