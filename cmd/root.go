/*
Copyright © 2024 Eric Yager
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sess",
	Short: "A command line tool for synchronizing development environments.",
	Long: `SESS (Synchronize Environments Simply, Stupid) is a command line tool
for synchronizing development environments using YAML configuration files.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(installCmd)
}
