package commands

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"

	"github.com/getstackhead/stackhead/cli/ansible"
	"github.com/getstackhead/stackhead/cli/routines"
)

var SetupServer = &cobra.Command{
	Use:     "setup [ipv4 address]",
	Example: "setup 192.168.178.14",
	Short:   "Prepare a server for deployment",
	Long:    `setup will install all required software on a server. You are then able to deploy projects onto it.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		routines.RunTask(
			routines.Text(fmt.Sprintf("Deploying to server at IP \"%s\"", args[0])),
			routines.Execute(func(wg *sync.WaitGroup, result chan routines.TaskResult) {
				defer wg.Done()

				// Generate Inventory file
				inventoryFile, err := ansible.CreateInventoryFile(ansible.IpAddress(args[0]))
				if err == nil {
					defer os.Remove(inventoryFile)
					err = routines.ExecAnsiblePlaybook("server-provision", inventoryFile)
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
