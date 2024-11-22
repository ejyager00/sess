/*
Copyright © 2024 Eric Yager
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
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

func versionSatisfiesConstraint(versionConstraint string, version string) (bool, error) {
	constraints, err := semver.NewConstraint(versionConstraint)
	if err != nil {
		return false, fmt.Errorf("invalid version constraint: %w", err)
	}

	re := regexp.MustCompile(`(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?`)
	match := re.FindString(version)
	if match == "" {
		return false, fmt.Errorf("no semantic version found in string: %s", version)
	}

	versionObj, err := semver.NewVersion(match)
	if err != nil {
		return false, fmt.Errorf("error parsing version: %v", err)
	}

	return constraints.Check(versionObj), nil
}

func checkVersion(version_check_command string, expected_version string) (bool, error) {
	output, err := exec.Command("sh", "-c", version_check_command).Output()
	if err != nil {
		return false, fmt.Errorf("error running version check command: %v", err)
	}

	actual_version := strings.TrimSpace(string(output))

	return versionSatisfiesConstraint(expected_version, actual_version)
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
