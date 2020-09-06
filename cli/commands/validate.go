package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/getstackhead/stackhead/cli/ansible"
	"github.com/getstackhead/stackhead/cli/jsonschema"
)

// Validate is a command object for Cobra that provides the validate command
var Validate = &cobra.Command{
	Use:     "validate [path to project definition file]",
	Example: "validate ./my-project-definition.yml",
	Short:   "Validate a project definition file",
	Long:    `validate is used to make sure your project definition file meets the StackHead project definition syntax.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var collectionDir, err = ansible.GetStackHeadCollectionLocation()
		if err != nil {
			panic(err.Error())
		}
		result, err := jsonschema.ValidateFile(collectionDir, args[0])

		if err != nil {
			panic(err.Error())
		}

		errorMessage := jsonschema.ShouldValidate(result)
		if len(errorMessage) == 0 {
			_, err = fmt.Fprintln(os.Stdout, "The project definition is valid")
		} else {
			_, err = fmt.Fprintln(os.Stderr, errorMessage)
			os.Exit(1)
		}
		if err != nil {
			panic(err.Error())
		}
	},
}
