package project

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/getstackhead/stackhead/ansible"
	"github.com/getstackhead/stackhead/routines"
)

// DeployApplication is a command object for Cobra that provides the deploy command
var DeployApplication = &cobra.Command{
	Use:     "deploy [path to project definition] [ipv4 address]",
	Example: "deploy ./my_project.yml 192.168.178.14",
	Short:   "Deploy a project onto the target server",
	Long:    `deploy will deploy the given project onto the server`,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		routines.RunTask(routines.Task{
			Name: fmt.Sprintf("Deploying project \"%s\" onto server with IP \"%s\"", args[0], args[1]),
			Run: func(r routines.RunningTask) error {
				// Generate Inventory file
				inventoryFile, err := ansible.CreateInventoryFile(args[1], args[0])
				if err == nil {
					defer os.Remove(inventoryFile)
					err = routines.ExecAnsiblePlaybook("application-deploy", inventoryFile, nil)
				}

				return err
			},
		})
	},
}
