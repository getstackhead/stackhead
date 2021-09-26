package project

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/getstackhead/stackhead/routines"
)

// Validate is a command object for Cobra that provides the validate command
func Validate() *cobra.Command {
	var version, branch string
	var ignoreSslCertificate bool
	var command = &cobra.Command{
		Use:     "validate [path to project definition file]",
		Example: "validate ./my-project-definition.yml",
		Short:   "Validate a project definition file",
		Long:    `validate is used to make sure your project definition file meets the StackHead project definition syntax.`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			routines.CobraValidationBase(
				"project-definition.schema.json",
				version,
				branch,
				ignoreSslCertificate,
			)(cmd, args)

			if !strings.HasSuffix(args[0], ".stackhead.yml") && !strings.HasSuffix(args[0], ".stackhead.yaml") {
				panic("The file name must end in \".stackhead.yml\" or \".stackhead.yaml\"!")
			}
		},
	}
	command.PersistentFlags().StringVar(&version, "version", "", "Version of schema to use (requires internet connection)")
	command.PersistentFlags().StringVar(&branch, "branch", "", "Branch of schema to use (requires internet connection)")
	command.PersistentFlags().BoolVar(&ignoreSslCertificate, "ignore-ssl-certificate", false, "Whether to ignore the SSL certificate for Web request (when --version) is used")

	return command
}
