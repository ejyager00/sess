/*
Copyright © 2024 Eric Yager
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/ejyager00/sess/internal/io"
	"github.com/ejyager00/sess/internal/models"
	"github.com/spf13/cobra"
)

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

func checkVersion(versionCheckCommand string, expectedVersion string) (bool, error) {
	output, err := exec.Command("sh", "-c", versionCheckCommand).Output()
	if err != nil {
		return false, fmt.Errorf("error running version check command: %v", err)
	}

	actualVersion := strings.TrimSpace(string(output))
	return versionSatisfiesConstraint(expectedVersion, actualVersion)
}

func validateTools(tools []models.Tool) {
	for _, tool := range tools {
		matches, err := checkVersion(tool.VersionCheckCommand, tool.Version)
		if err != nil {
			fmt.Printf("%s: error checking version: %v\n", tool.Name, err)
		} else {
			fmt.Printf("%s: %v\n", tool.Name, matches)
		}
	}
}

func validateEnvVariables(envVariables map[string]string) {
	for envVar, expectedValue := range envVariables {
		actualValue, exists := os.LookupEnv(envVar)
		if !exists {
			fmt.Printf("Environment variable %s is not set\n", envVar)
		} else if actualValue != expectedValue {
			fmt.Printf("Environment variable %s has incorrect value. Expected: %s, Got: %s\n",
				envVar, expectedValue, actualValue)
		} else {
			fmt.Printf("Environment variable %s is correctly set\n", envVar)
		}
	}
}

func validateDotfiles(dotfiles []string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting user home directory: %v", err)
	}

	for _, dotfile := range dotfiles {
		dotfilePath := filepath.Join(homeDir, dotfile)
		if _, err := os.Stat(dotfilePath); err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("Dotfile %s does not exist\n", dotfile)
			} else {
				fmt.Printf("Error checking dotfile %s: %v\n", dotfile, err)
			}
		} else {
			fmt.Printf("Dotfile %s exists\n", dotfile)
		}
	}
	return nil
}

func validateExtensions(extensions []models.Extension) error {
	// Get installed VS Code extensions
	cmd := exec.Command("codium", "--list-extensions")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error getting VS Code extensions: %v", err)
	}

	// Parse installed extensions into a map for easy lookup
	installedExtensions := make(map[string]bool)
	for _, ext := range strings.Split(string(output), "\n") {
		ext = strings.TrimSpace(ext)
		if ext != "" {
			installedExtensions[strings.ToLower(ext)] = true
		}
	}

	// Check each required extension
	for _, extension := range extensions {
		if extension.Ide != "vscode" {
			fmt.Printf("Skipping extension %s - only VS Code extensions supported currently\n", extension.Id)
			continue
		}

		if installedExtensions[strings.ToLower(extension.Id)] {
			fmt.Printf("Extension %s is installed\n", extension.Id)
		} else {
			fmt.Printf("Extension %s is not installed\n", extension.Id)
		}
	}
	return nil
}

var validateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Validate the current environment against the configuration.",
	Long:  `Validate that the current environment matches the provided configuration file.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		schema, err := io.ParseYamlFile(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		validateTools(schema.Tools)
		validateEnvVariables(schema.EnvVariables)

		if err := validateDotfiles(schema.Dotfiles); err != nil {
			fmt.Println(err)
		}

		if err := validateExtensions(schema.Extensions); err != nil {
			fmt.Println(err)
		}
	},
}
