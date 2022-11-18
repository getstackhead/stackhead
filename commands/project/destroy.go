package project

import (
	"fmt"

	xfs "github.com/saitho/golang-extended-fs/v2"
	"github.com/spf13/cobra"

	"github.com/getstackhead/stackhead/commands"
	"github.com/getstackhead/stackhead/project"
	"github.com/getstackhead/stackhead/routines"
	"github.com/getstackhead/stackhead/system"
	"github.com/getstackhead/stackhead/terraform"
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

		routines.RunTask(routines.Task{
			Name: fmt.Sprintf("Destroying project \"%s\" on server with IP \"%s\"", args[0], args[1]),
			Run: func(r routines.RunningTask) error {
				// Init modules
				for _, module := range system.Context.GetModulesInOrder() {
					moduleSettings := system.GetModuleSettings(module.GetConfig().Name)
					module.Init(moduleSettings)
				}

				hasTerraformDir, _ := xfs.HasFolder("ssh://" + projectDefinition.GetTerraformDirectoryPath())
				if hasTerraformDir {
					if err := routines.RunTask(routines.Task{
						Name: "Destroying Terraform plans",
						Run: func(r routines.RunningTask) error {
							if err := terraform.Init(projectDefinition.GetTerraformDirectoryPath()); err != nil {
								return err
							}
							if err := terraform.Destroy(projectDefinition.GetTerraformDirectoryPath()); err != nil {
								return err
							}
							return nil
						},
					}); err != nil {
						return err
					}
				}

				hasProjectDir, _ := xfs.HasFolder("ssh://" + projectDefinition.GetDirectoryPath())
				if hasProjectDir {
					// Removing project directory
					if err := xfs.DeleteFolder("ssh://"+projectDefinition.GetDirectoryPath(), true); err != nil {
						return err
					}
				}

				// Run destroy scripts from DNS plugins
				for _, module := range system.Context.DNSModules {
					if module.GetConfig().Type != "dns" {
						continue
					}
					moduleSettings := system.GetModuleSettings(module.GetConfig().Name)
					if err := module.(system.DNSModule).Destroy(moduleSettings); err != nil {
						return err
					}
				}

				return nil
			},
		})
	},
}
