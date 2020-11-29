package cli

import (
	"github.com/spf13/cobra"
)

func GetCommands() *cobra.Command {
	command := &cobra.Command{
		Use:     "cli",
		Short:   "StackHead CLI commands",
	}
	command.AddCommand(Validate)
	return command
}
