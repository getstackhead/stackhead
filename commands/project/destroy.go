package project

import (
	"fmt"

	xfs "github.com/saitho/golang-extended-fs/v2"
	"github.com/spf13/cobra"

	"github.com/getstackhead/stackhead/commands"
	"github.com/getstackhead/stackhead/project"
	"github.com/getstackhead/stackhead/routines"
	"github.com/getstackhead/stackhead/system"
)

// DestroyApplication is a command object for Cobra that provides the destroy command
var DestroyApplication = &cobra.Command{
	Use:     "destroy [path to project definition] [ipv4 address]",
	Example: "destroy ./my_project.yml 192.168.178.14",
	Short:   "Destroy a deployed project on a target server",
	Long:    `destroy will tear down the given project that has been deployed onto the server`,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		projectDefinition, err := project.LoadProjectDefinition(args[0])
		if err != nil {
			panic("unable to load project definition file. " + err.Error())
		}
		commands.PrepareContext(args[1], system.ContextActionProjectDeploy, projectDefinition)

		modules := system.Context.GetModulesInOrder()
		for i, j := 0, len(modules)-1; i < j; i, j = i+1, j-1 { // reverse module list
			modules[i], modules[j] = modules[j], modules[i]
		}

		// Init modules
		for _, module := range modules {
			moduleSettings := system.GetModuleSettings(module.GetConfig().Name)
			module.Init(moduleSettings)
		}
		taskRunner := routines.TaskRunner{}

		subTasks := []routines.Task{}

		if hasProjectDir, _ := xfs.HasFolder("ssh://" + projectDefinition.GetDirectoryPath()); hasProjectDir {

			// Run destroy scripts from plugins
			for _, module := range modules {
				moduleSettings := system.GetModuleSettings(module.GetConfig().Name)
				subTasks = append(subTasks, routines.Task{
					Name: "Remove module configurations for " + module.GetConfig().Name,
					Run: func(r *routines.Task) error {
						return module.Destroy(moduleSettings)
					},
					IsSubtask:           true,
					ErrorAsErrorMessage: true,
				})
			}

			subTasks = append(subTasks, routines.Task{
				Name: "Removing project directory",
				Run: func(r *routines.Task) error {
					return xfs.DeleteFolder("ssh://"+projectDefinition.GetDirectoryPath(), true)
				},
				IsSubtask:           true,
				ErrorAsErrorMessage: true,
			})
		}

		_ = taskRunner.RunTask(routines.Task{
			Name: fmt.Sprintf("Destroying project \"%s\" on server with IP \"%s\"", args[0], args[1]),
			Run: func(r *routines.Task) error {
				return nil
			},
			SubTasks: subTasks,
			//RunAllSubTasksDespiteError: true,
		})
	},
}
