package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/getstackhead/stackhead/cli/commands"
)

var cfgFile string
var collectionPath string
var verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "stackhead-cli",
	Short: "A brief description of your application",
	Long: `Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

// main adds all child commands to the root command and sets flags appropriately.
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is ./.stackhead-cli.yaml or $HOME/.stackhead-cli.yaml)")
	rootCmd.PersistentFlags().StringVarP(&collectionPath, "collection-path", "p", "", "Path to Ansible collection directory")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Show more output")

	rootCmd.AddCommand(commands.Validate)
	rootCmd.AddCommand(commands.Init())
	rootCmd.AddCommand(commands.SetupServer)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if collectionPath != "" {
		viper.Set("ansible.collection_path", collectionPath)
	}
	if verbose {
		viper.Set("verbose", verbose)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.SetConfigName(".stackhead-cli") // .stackhead-cli.yml (no file extension needed)
		// Search config in working directory
		if dir, err := os.Getwd(); err == nil {
			viper.AddConfigPath(dir)
		}
		// Search config in current or home directory
		viper.AddConfigPath(home)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
