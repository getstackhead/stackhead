package project

import (
	"github.com/spf13/cobra"

	"github.com/getstackhead/stackhead/cli/routines"
)

// Validate is a command object for Cobra that provides the validate command
var Validate = &cobra.Command{
	Use:     "validate [path to project definition file]",
	Example: "validate ./my-project-definition.yml",
	Short:   "Validate a project definition file",
	Long:    `validate is used to make sure your project definition file meets the StackHead project definition syntax.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		routines.Validate(args[0], "project-definition.schema.json")
	},
}
