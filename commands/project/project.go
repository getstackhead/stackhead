package project

import (
	"embed"

	"github.com/spf13/cobra"
)

func GetCommands(LocalSchemas embed.FS) *cobra.Command {
	command := &cobra.Command{
		Use:   "project",
		Short: "Project commands",
	}
	command.AddCommand(DeployApplication, DestroyApplication, Validate(LocalSchemas)) // nolint:typecheck
	return command
}
