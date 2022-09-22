package project

import (
	"fmt"

	xfs "github.com/saitho/golang-extended-fs/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/getstackhead/stackhead/commands"
	"github.com/getstackhead/stackhead/project"
	"github.com/getstackhead/stackhead/routines"
	"github.com/getstackhead/stackhead/system"
	"github.com/getstackhead/stackhead/terraform"
)

// DeployApplication is a command object for Cobra that provides the deploy command
var DeployApplication = &cobra.Command{
	Use:     "deploy [path to project definition] [ipv4 address]",
	Example: "deploy ./my_project.yml 192.168.178.14",
	Short:   "Deploy a project onto the target server",
	Long:    `deploy will deploy the given project onto the server`,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		projectDefinition, err := project.LoadProjectDefinition(args[0])
		if err != nil {
			panic("unable to load project definition file. " + err.Error())
		}
		commands.PrepareContext(args[1], system.ContextActionProjectDeploy, projectDefinition)
		_ = routines.RunTask(routines.Task{
			Name: fmt.Sprintf("Deploying project \"%s\" onto server with IP \"%s\"", args[0], args[1]),
			Run: func(r routines.RunningTask) error {
				// Validate StackHead version
				isValid, err := system.ValidateVersion()
				if err != nil {
					logger.Debugln(err)
					return fmt.Errorf("Unable to validate StackHead version.")
				}
				if !isValid {
					return fmt.Errorf("Trying to deploy with a newer version of StackHead than used for server setup. Please run a server setup again.")
				}

				// Init modules
				for _, module := range system.Context.GetModulesInOrder() {
					moduleSettings := system.GetModuleSettings(module.GetConfig().Name)
					module.Init(moduleSettings)
				}

				if err := routines.RunTask(routines.Task{
					Name: "Preparing project structure",
					Run: func(r routines.RunningTask) error {
						var err error

						r.PrintLn("Create project directory if not exists")
						if err := xfs.CreateFolder("ssh://" + projectDefinition.GetDirectoryPath()); err != nil {
							return err
						}

						r.PrintLn("Prepare Terraform directory")
						if err := xfs.CreateFolder("ssh://" + projectDefinition.GetTerraformDirectoryPath()); err != nil {
							return err
						}

						return err
					},
					ErrorAsErrorMessage: true,
				}); err != nil {
					return err
				}

				if err := routines.RunTask(routines.Task{
					Name: "Generating Terraform files",
					Run: func(r routines.RunningTask) error {
						// Collect exposed services
						var exposedServices []project.DomainExpose
						for _, domain := range projectDefinition.Domains {
							exposedServices = append(exposedServices, domain.Expose...)
						}

						r.PrintLn("Symlinking Terraform providers")
						if err := terraform.SymlinkProviders(system.Context.Project); err != nil {
							return fmt.Errorf("Unable to symlink Terraform providers")
						}

						r.PrintLn("Generate project-specific Terraform providers")
						modulesWithProviders := terraform.FilterModulesWithProviders(system.Context.GetModulesInOrder())
						fileContents, err := terraform.BuildProviders(modulesWithProviders, terraform.ONLY_PER_PROJECT)
						if err != nil {
							return fmt.Errorf("Unable to generate project-specific Terraform providers: " + err.Error())
						}
						if fileContents.Len() > 0 {
							if err := xfs.WriteFile("ssh://"+projectDefinition.GetTerraformProjectProvidersFilePath(), fileContents.String()); err != nil {
								return err
							}
						}

						for _, module := range system.Context.GetModulesInOrder() {
							if module.GetConfig().Type == "plugin" {
								continue
							}
							moduleSettings := system.GetModuleSettings(module.GetConfig().Name)
							if err := module.Deploy(moduleSettings); err != nil {
								return err
							}
						}
						return nil
					},
				}); err != nil {
					return err
				}

				// todo: integrate stackhead_config.terraform.manual_update
				if err := routines.RunTask(routines.Task{
					Name: "Applying Terraform plans",
					Run: func(r routines.RunningTask) error {
						if err := terraform.Init(system.Context.Project.GetTerraformDirectoryPath()); err != nil {
							return err
						}
						if err := terraform.Apply(system.Context.Project.GetTerraformDirectoryPath()); err != nil {
							return err
						}
						return nil
					},
				}); err != nil {
					return err
				}

				/**
				- name: Create project specific crontab
					include_tasks: "../roles/config_terraform/tasks/setup_crontab.yml"
					when: ensure == 'present'
				- name: Remove project specific crontab
					include_tasks: "../roles/config_terraform/tasks/remove_crontab.yml"
					when: ensure == 'absent'
				*/

				return nil
			},
		})
	},
}
