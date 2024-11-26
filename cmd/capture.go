/*
Copyright © 2024 Eric Yager
*/
package cmd

import "github.com/spf13/cobra"

var captureCmd = &cobra.Command{
	Use:   "capture",
	Short: "Capture the current environment.",
	Long:  `Captures the current relevant environment state and stores it in a configuration YAML file.`,
}
