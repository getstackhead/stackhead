package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/getstackhead/stackhead/ansible"
	"github.com/getstackhead/stackhead/routines"
)

// DestroyApplication is a command object for Cobra that provides the destroy command
var DestroyApplication = &cobra.Command{
	Use:     "destroy [path to project definition] [ipv4 address]",
	Example: "destroy ./my_project.yml 192.168.178.14",
	Short:   "Destroy a deployed project on a target server",
	Long:    `destroy will tear down the given project that has been deployed onto the server`,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		routines.RunTask(routines.Task{
			Name: fmt.Sprintf("Destroying project \"%s\" on server with IP \"%s\"", args[0], args[1]),
			Run: func(r routines.RunningTask) error {
				// Generate Inventory file
				inventoryFile, err := ansible.CreateInventoryFile(args[1], args[0])

				if err == nil {
					defer os.Remove(inventoryFile)
					options := make(map[string]string)
					options["project_name"] = strings.TrimSuffix(strings.TrimSuffix(filepath.Base(args[0]), ".stackhead.yml"), ".stackhead.yaml")
					err = routines.ExecAnsiblePlaybook("application-destroy", inventoryFile, options)
				}
				return err
			},
		})
	},
}
