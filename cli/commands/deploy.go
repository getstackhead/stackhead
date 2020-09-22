package commands

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"

	"github.com/getstackhead/stackhead/cli/ansible"
	"github.com/getstackhead/stackhead/cli/routines"
)

// DeployApplication is a command object for Cobra that provides the deploy command
var DeployApplication = &cobra.Command{
	Use:     "deploy [path to project definition] [ipv4 address]",
	Example: "deploy ./my_project.yml 192.168.178.14",
	Short:   "Deploy a project onto the target server",
	Long:    `deploy will deploy the given project onto the server`,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		routines.RunTask(
			routines.Text(fmt.Sprintf("Deploying project \"%s\" onto server with IP \"%s\"", args[0], args[1])),
			routines.Execute(func(wg *sync.WaitGroup, result chan routines.TaskResult) {
				defer wg.Done()

				// Generate Inventory file
				inventoryFile, err := ansible.CreateInventoryFile(
					ansible.IPAddress(args[1]),
					ansible.ProjectDefinitionFile(args[0]),
				)
				if err == nil {
					defer os.Remove(inventoryFile)
					err = routines.ExecAnsiblePlaybook("application-deploy", inventoryFile, nil)
				}

				taskResult := routines.TaskResult{
					Error: err != nil,
				}
				if err != nil {
					taskResult.Message = err.Error()
				}

				result <- taskResult
			}),
		)
	},
}
