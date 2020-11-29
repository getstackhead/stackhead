package cli

import (
	"github.com/spf13/cobra"

	"github.com/getstackhead/stackhead/cli/routines"
)

// Validate is a command object for Cobra that provides the validate command
var Validate = &cobra.Command{
	Use:     "validate [path to StackHead CLI configuration file]",
	Example: "validate ./stackhead-module.yml",
	Short:   "Validate a StackHead module file",
	Long:    `validate is used to make sure your StackHead CLI configuration file meets the required syntax.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		routines.Validate(args[0], "cli-config.schema.json")
	},
}
