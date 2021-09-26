package module

import (
	"github.com/getstackhead/stackhead/routines"
	"github.com/spf13/cobra"
)

// Validate is a command object for Cobra that provides the validate command
func Validate() *cobra.Command {
	var version, branch string
	var ignoreSslCertificate bool
	var command = &cobra.Command{
		Use:     "validate [path to StackHead module file]",
		Example: "validate ./stackhead-module.yml",
		Short:   "Validate a StackHead module file",
		Long:    `validate is used to make sure your StackHead module file meets the required syntax.`,
		Args:    cobra.ExactArgs(1),
		Run: routines.CobraValidationBase(
			"module-config.schema.json",
			version,
			branch,
			ignoreSslCertificate,
		),
	}
	command.PersistentFlags().StringVar(&version, "version", "", "Version of schema to use (requires internet connection)")
	command.PersistentFlags().StringVar(&branch, "branch", "", "Branch of schema to use (requires internet connection)")
	command.PersistentFlags().BoolVar(&ignoreSslCertificate, "ignore-ssl-certificate", false, "Whether to ignore the SSL certificate for Web request (when --version) is used")

	return command
}
