package commands

import (
	"github.com/spf13/cobra"

	commandsinit "github.com/getstackhead/stackhead/cli/commands/init"
	"github.com/getstackhead/stackhead/cli/routines"
)

// Init is a command object for Cobra that provides the init command
func Init() *cobra.Command {
	version := ""
	command := &cobra.Command{
		Use:   "init",
		Short: "Install StackHead dependencies according to configuration file",
		Long: `init will install all required dependencies according to configuration file.
If no configuration file exists, it will start a wizard to create one.`,
		Run: func(cmd *cobra.Command, args []string) {
			routines.RunTask(commandsinit.InstallCollection(version)...)
			routines.RunTask(commandsinit.InstallModules...)
		},
	}
	command.PersistentFlags().StringVar(&version, "version", "", "Version of StackHead to be installed")
	return command
}
