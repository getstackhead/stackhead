package project

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/getstackhead/stackhead/config"
	"github.com/getstackhead/stackhead/plugins"
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
		config, err := config.LoadProjectDefinition(args[0])
		if err != nil {
			panic("unable to load project definition file. " + err.Error())
		}
		system.InitializeContext(args[1], system.ContextActionProjectDeploy, config)
		routines.RunTask(routines.Task{
			Name: fmt.Sprintf("Deploying project \"%s\" onto server with IP \"%s\"", args[0], args[1]),
			Run: func(r routines.RunningTask) error {
				var err error

				p, err := plugins.LoadPlugins()
				if err != nil {
					return err
				}
				for _, plugin := range p {
					if plugin.DeployProgram != nil {
						if err := plugin.DeployProgram.Run(nil); err != nil {
							fmt.Println(err)
							os.Exit(1)
						}
					}
				}

				return err
			},
		})
	},
}
