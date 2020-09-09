package commands

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"

	"github.com/getstackhead/stackhead/cli/ansible"
	"github.com/getstackhead/stackhead/cli/routines"
)

// DestroyApplication is a command object for Cobra that provides the destroy command
var DestroyApplication = &cobra.Command{
	Use:     "destroy [project name] [ipv4 address]",
	Example: "destroy my_project_name 192.168.178.14",
	Short:   "Destroy a deployed project on a target server",
	Long:    `destroy will tear down the given project that has been deployed onto the server`,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		routines.RunTask(
			routines.Text(fmt.Sprintf("Destroying project \"%s\" on server with IP \"%s\"", args[0], args[1])),
			routines.Execute(func(wg *sync.WaitGroup, result chan routines.TaskResult) {
				defer wg.Done()

				// Generate Inventory file
				inventoryFile, err := ansible.CreateInventoryFile(
					ansible.IPAddress(args[1]),
				)
				if err == nil {
					defer os.Remove(inventoryFile)
					options := make(map[string]string)
					options["project_name"] = args[0]
					err = routines.ExecAnsiblePlaybook("application-destroy", inventoryFile, options)
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
