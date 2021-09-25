package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/getstackhead/stackhead/ansible"
	"github.com/getstackhead/stackhead/routines"
)

// SetupServer is a command object for Cobra that provides the setup command
var SetupServer = &cobra.Command{
	Use:     "setup [ipv4 address]",
	Example: "setup 192.168.178.14",
	Short:   "Prepare a server for deployment",
	Long:    `setup will install all required software on a server. You are then able to deploy projects onto it.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		routines.RunTask(routines.Task{
			Name: fmt.Sprintf("Setting up server at IP \"%s\"", args[0]),
			Run: func(r routines.RunningTask) error {
				// Generate Inventory file
				inventoryFile, err := ansible.CreateInventoryFile(args[0], "")
				if err == nil {
					defer os.Remove(inventoryFile)
					err = routines.ExecAnsiblePlaybook("server-provision", inventoryFile, nil)
				}
				return err
			},
			ErrorAsErrorMessage: true,
		})
	},
}
