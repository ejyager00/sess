package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func parseYamlFile(filePath string) (interface{}, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var parsedYaml interface{}
	err = yaml.Unmarshal(data, &parsedYaml)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML: %v", err)
	}

	return parsedYaml, nil
}

var validateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Validate the current environment against the configuration.",
	Long:  `Validate that the current environment matches the provided configuration file.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		parsedYaml, err := parseYamlFile(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		output, err := yaml.Marshal(parsedYaml)
		if err != nil {
			fmt.Printf("Error formatting YAML: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(output))
	},
}
