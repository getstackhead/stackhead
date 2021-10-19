package project

import (
	"fmt"
	"github.com/getstackhead/stackhead/config"
	"github.com/getstackhead/stackhead/stackhead"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

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
		config, err := config.LoadProjectDefinition(args[0])
		if err != nil {
			panic("unable to load project definition file. " + err.Error())
		}
		stackhead.InitializeContext(args[1], stackhead.ContextActionProjectDestroy, config)
		routines.RunTask(routines.Task{
			Name: fmt.Sprintf("Destroying project \"%s\" on server with IP \"%s\"", args[0], args[1]),
			Run: func(r routines.RunningTask) error {
				var err error

				options := make(map[string]string)
				options["project_name"] = strings.TrimSuffix(strings.TrimSuffix(filepath.Base(args[0]), ".stackhead.yml"), ".stackhead.yaml")

				// todo: run destroy
				r.PrintLn("Destroy not yet implemented.")

				return err
			},
		})
	},
}
