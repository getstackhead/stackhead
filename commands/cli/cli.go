package cli

import (
	"embed"

	"github.com/spf13/cobra"
)

func GetCommands(LocalSchemas embed.FS) *cobra.Command {
	command := &cobra.Command{
		Use:   "cli",
		Short: "StackHead CLI commands",
	}
	validate := Validate(LocalSchemas) // nolint:typecheck
	command.AddCommand(validate)
	return command
}
