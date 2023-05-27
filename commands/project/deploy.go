package project

import (
	"fmt"
	"golang.org/x/exp/slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/getstackhead/stackhead/commands"
	"github.com/getstackhead/stackhead/project"
	"github.com/getstackhead/stackhead/routines"
	"github.com/getstackhead/stackhead/system"
)

// DeployApplication is a command object for Cobra that provides the deploy command
var DeployApplication = &cobra.Command{
	Use:     "deploy [path to project definition] [ipv4 address]",
	Example: "deploy ./my_project.yml 192.168.178.14",
	Short:   "Deploy a project onto the target server",
	Long:    `deploy will deploy the given project onto the server`,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		projectDefinition, err := project.LoadProjectDefinition(args[0])
		if err != nil {
			panic("unable to load project definition file. " + err.Error())
		}
		commands.PrepareContext(args[1], system.ContextActionProjectDeploy, projectDefinition)

		if err := routines.RunTask(routines.ValidateStackHeadVersionTask); err != nil {
			return
		}

		// Init modules
		for _, module := range system.Context.GetModulesInOrder() {
			moduleSettings := system.GetModuleSettings(module.GetConfig().Name)
			module.Init(moduleSettings)
		}

		err = routines.RunTask(routines.PrepareProjectTask(projectDefinition))
		if err != nil {
			return
		}
		err = routines.RunTask(routines.CollectResourcesTask(projectDefinition))
		if err != nil {
			return
		}

		// Confirm resource creation
		fmt.Println("\nStackHead will try to create or update the following resources:")
		for _, resource := range system.Context.Resources {
			operation := system.GetOperationLabel(resource)
			fmt.Println(fmt.Sprintf("- [%s] %s %s", operation, resource.Type, resource.Name))
		}
		fmt.Println("")
		fmt.Print("Please confirm with \"y\" or \"yes\": ")
		if askForConfirmation() {
			_ = routines.RunTask(routines.Task{
				Name: "Creating resources",
				Run: func(r routines.RunningTask) error {
					success, errors := system.ApplyResourcesFromContext()
					if success {
						return nil
					}
					errorMessages := []string{"The following errors occurred:"}
					for _, err2 := range errors {
						errorMessages = append(errorMessages, "- "+err2.Error())
					}
					return fmt.Errorf(strings.Join(errorMessages, "\n"))
				},
				ErrorAsErrorMessage: true,
			})
		}
	},
}

func askForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		fmt.Println(err)
	}
	if slices.Contains([]string{"y", "Y", "yes", "Yes", "YES"}, response) {
		return true
	} else if slices.Contains([]string{"n", "N", "no", "No", "NO"}, response) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return askForConfirmation()
	}
}
