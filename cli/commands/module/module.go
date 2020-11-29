package module

import (
	"github.com/spf13/cobra"
)

func GetCommands() *cobra.Command {
	command := &cobra.Command{
		Use:     "module",
		Short:   "StackHead module commands",
	}
	command.AddCommand(Validate)
	return command
}
