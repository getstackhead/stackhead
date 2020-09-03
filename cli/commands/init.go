package commands

import (
	"github.com/spf13/cobra"

	"github.com/getstackhead/stackhead/cli/commands/init"
	"github.com/getstackhead/stackhead/cli/routines"
)

var Init = &cobra.Command{
	Use:   "init",
	Short: "Install StackHead dependencies according to configuration file",
	Long: `init will install all required dependencies according to configuration file.
If no configuration file exists, it will start a wizard to create one.`,
	Run: func(cmd *cobra.Command, args []string) {

		routines.RunTask(commands_init.InstallCollection...)
		routines.RunTask(commands_init.InstallModules...)
	},
}
