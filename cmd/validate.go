/*
Copyright © 2024 Eric Yager
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ejyager00/sess/models"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

func parseYamlFile(filePath string) (*models.EnvironmentSchema, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var schema models.EnvironmentSchema
	err = yaml.Unmarshal(data, &schema)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML: %v", err)
	}

	return &schema, nil
}

func checkVersion(version_check_command string, expected_version string) (bool, error) {
	// Execute the version check command
	output, err := exec.Command("sh", "-c", version_check_command).Output()
	if err != nil {
		return false, fmt.Errorf("error running version check command: %v", err)
	}

	// Convert output to string and trim whitespace
	actual_version := strings.TrimSpace(string(output))

	// If expected version ends in .x, only compare major version
	if strings.HasSuffix(expected_version, ".x") {
		major_version := strings.TrimSuffix(expected_version, ".x")
		return strings.HasPrefix(actual_version, major_version), nil
	}

	// Otherwise do exact match
	return actual_version == expected_version, nil
}

var validateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Validate the current environment against the configuration.",
	Long:  `Validate that the current environment matches the provided configuration file.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		schema, err := parseYamlFile(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, tool := range schema.Tools {
			matches, err := checkVersion(tool.VersionCheckCommand, tool.Version)
			if err != nil {
				fmt.Printf("%s: error checking version: %v\n", tool.Name, err)
				continue
			}

			fmt.Printf("%s: %v\n", tool.Name, matches)
		}
	},
}
