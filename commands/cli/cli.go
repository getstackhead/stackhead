package cli

import (
	"github.com/spf13/cobra"
)

func GetCommands() *cobra.Command {
	command := &cobra.Command{
		Use:   "cli",
		Short: "StackHead CLI commands",
	}
	validate := Validate() // nolint:typecheck
	command.AddCommand(validate)
	return command
}
