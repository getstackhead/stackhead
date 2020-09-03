package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/getstackhead/stackhead/cli/jsonschema"
)

var Validate = &cobra.Command{
	Use:   "validate [path to project definition file]",
	Short: "Validate a project definition file",
	Long:  `validate is used to make sure your project definition file meets the StackHead project definition syntax.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := jsonschema.ValidateFile(args[0])

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
