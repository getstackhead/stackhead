package project

import (
	"github.com/spf13/cobra"
)

func GetCommands() *cobra.Command {
	command := &cobra.Command{
		Use:   "project",
		Short: "Project commands",
	}
	command.AddCommand(DeployApplication, DestroyApplication, Validate()) // nolint:typecheck
	return command
}
