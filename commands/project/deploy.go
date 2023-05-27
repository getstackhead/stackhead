package project

import (
	"fmt"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"

	"github.com/getstackhead/stackhead/commands"
	"github.com/getstackhead/stackhead/project"
	"github.com/getstackhead/stackhead/routines"
	"github.com/getstackhead/stackhead/system"
)

// DeployApplication is a command object for Cobra that provides the deploy command
var DeployApplication = func() *cobra.Command {
	var autoConfirm bool
	var noRollback bool
	command := &cobra.Command{
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
			taskRunner := routines.TaskRunner{}

			if err := taskRunner.RunTask(routines.ValidateStackHeadVersionTask); err != nil {
				return
			}

			// Init modules
			for _, module := range system.Context.GetModulesInOrder() {
				moduleSettings := system.GetModuleSettings(module.GetConfig().Name)
				module.Init(moduleSettings)
			}

			err = taskRunner.RunTask(routines.PrepareProjectTask(projectDefinition))
			if err != nil {
				return
			}
			err = taskRunner.RunTask(routines.CollectResourcesTask(projectDefinition))
			if err != nil {
				return
			}

			if autoConfirm {
				_ = taskRunner.RunTask(routines.CreateResources)
			} else {
				// Confirm resource creation
				fmt.Println("\nStackHead will try to create or update the following resources:")
				for _, resourceGroup := range system.Context.Resources {
					for _, resource := range resourceGroup.Resources {
						fmt.Println(fmt.Sprintf("- %s", resource.ToString(false)))
					}
				}
				fmt.Println("")
				fmt.Print("Please confirm with \"y\" or \"yes\": ")
				if askForConfirmation() {
					_ = taskRunner.RunTask(routines.CreateResources)
				}
			}

			if !noRollback {
				// Rollback may be skipped if CreateResources does not trigger a rollback
				_ = taskRunner.RunTask(routines.RollbackResources)
			}
		},
	}
	command.PersistentFlags().BoolVar(&autoConfirm, "autoconfirm", false, "Whether to auto-confirm resource changes")
	command.PersistentFlags().BoolVar(&noRollback, "no-rollback", false, "Do not rollback on errors")
	return command
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
