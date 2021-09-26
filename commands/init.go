package commands

import (
	commandsinit "github.com/getstackhead/stackhead/commands/init"
	"github.com/getstackhead/stackhead/routines"
	"github.com/spf13/cobra"
)

// Init is a command object for Cobra that provides the init command
func Init() *cobra.Command {
	command := &cobra.Command{
		Use:   "init",
		Short: "Install StackHead dependencies according to configuration file",
		Long: `init will install all required dependencies according to configuration file.
If no configuration file exists, it will start a wizard to create one.`,
		Run: func(cmd *cobra.Command, args []string) {
			routines.RunTask(commandsinit.InstallPlugins())
		},
	}
	return command
}
