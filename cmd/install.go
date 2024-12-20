/*
Copyright © 2024 Eric Yager
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/ejyager00/sess/internal/io"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [file]",
	Short: "Install the tools and extensions specified in the configuration file.",
	Long:  `Install the tools and extensions specified in the configuration file.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		schema, err := io.ParseYamlFile(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for key, value := range schema.EnvVariables {
			err := os.Setenv(key, value)
			if err != nil {
				fmt.Printf("Error setting environment variable %s: %v\n", key, err)
				os.Exit(1)
			}
		}
	},
}
