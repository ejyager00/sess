/*
Copyright © 2024 Eric Yager
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/ejyager00/sess/internal/models"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

// listDotFiles lists all dot files in the given directory and returns a slice of their names
func listDotFiles(dir string) ([]string, error) {
	var dotFiles []string

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}

	for _, file := range files {
		name := file.Name()
		if strings.HasPrefix(name, ".") {
			dotFiles = append(dotFiles, name)
		}
	}

	return dotFiles, nil
}

// promptForDotFileSelection prints the list of dot files and gets user selection
func promptForDotFileSelection(dotFiles []string) ([]string, error) {
	if len(dotFiles) == 0 {
		return nil, fmt.Errorf("no dot files found in directory")
	}

	fmt.Println("Found the following dot files:")
	for i, file := range dotFiles {
		fmt.Printf("%d. %s\n", i+1, file)
	}

	fmt.Print("\nEnter comma-separated numbers of files to include (e.g. 1,3,4): ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading input: %v", err)
	}

	// Parse selection
	selections := strings.Split(strings.TrimSpace(input), ",")
	var selectedFiles []string

	for _, sel := range selections {
		num, err := strconv.Atoi(strings.TrimSpace(sel))
		if err != nil {
			return nil, fmt.Errorf("invalid selection: %s", sel)
		}
		if num < 1 || num > len(dotFiles) {
			return nil, fmt.Errorf("selection out of range: %d", num)
		}
		selectedFiles = append(selectedFiles, dotFiles[num-1])
	}

	return selectedFiles, nil
}

// captureDotFiles handles the dotfile capture workflow
func captureDotFiles(dir string) ([]string, error) {
	dotFiles, err := listDotFiles(dir)
	if err != nil {
		return nil, err
	}

	selectedFiles, err := promptForDotFileSelection(dotFiles)
	if err != nil {
		return nil, err
	}

	fmt.Println("\nSelected dot files:")
	for _, file := range selectedFiles {
		fmt.Println(file)
	}

	return selectedFiles, nil
}

// listEnvironmentVariables lists all environment variables and returns them as a map
func listEnvironmentVariables() map[string]string {
	envVars := make(map[string]string)
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			envVars[pair[0]] = pair[1]
		}
	}
	return envVars
}

// promptForEnvVarSelection prints the list of environment variables and gets user selection
func promptForEnvVarSelection(envVars map[string]string) (map[string]string, error) {
	if len(envVars) == 0 {
		return nil, fmt.Errorf("no environment variables found")
	}

	// Convert map to sorted slice for consistent display
	var envVarList []string
	for key := range envVars {
		envVarList = append(envVarList, key)
	}
	sort.Strings(envVarList)

	fmt.Println("\nFound the following environment variables:")
	for i, key := range envVarList {
		fmt.Printf("%d. %s=%s\n", i+1, key, envVars[key])
	}

	fmt.Print("\nEnter comma-separated numbers of environment variables to include (e.g. 1,3,4): ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading input: %v", err)
	}

	// Parse selection
	selections := strings.Split(strings.TrimSpace(input), ",")
	selectedVars := make(map[string]string)

	for _, sel := range selections {
		num, err := strconv.Atoi(strings.TrimSpace(sel))
		if err != nil {
			return nil, fmt.Errorf("invalid selection: %s", sel)
		}
		if num < 1 || num > len(envVarList) {
			return nil, fmt.Errorf("selection out of range: %d", num)
		}
		key := envVarList[num-1]
		selectedVars[key] = envVars[key]
	}

	return selectedVars, nil
}

// captureEnvVariables handles the environment variable capture workflow
func captureEnvVariables() (map[string]string, error) {
	envVars := listEnvironmentVariables()
	selectedVars, err := promptForEnvVarSelection(envVars)
	if err != nil {
		return nil, err
	}

	fmt.Println("\nSelected environment variables:")
	for key, value := range selectedVars {
		fmt.Printf("%s=%s\n", key, value)
	}

	return selectedVars, nil
}

// buildEnvironmentSchema constructs the environment schema from captured data
func buildEnvironmentSchema(dotfiles []string, envVars map[string]string) *models.EnvironmentSchema {
	return &models.EnvironmentSchema{
		Dotfiles:     dotfiles,
		EnvVariables: envVars,
		// Future fields will be added here:
		// Tools: capturedTools,
		// Extensions: capturedExtensions,
	}
}

// saveEnvironmentSchema saves the schema to a YAML file
func saveEnvironmentSchema(schema *models.EnvironmentSchema, filename string) error {
	yamlData, err := yaml.Marshal(schema)
	if err != nil {
		return fmt.Errorf("error marshaling YAML: %v", err)
	}

	err = os.WriteFile(filename, yamlData, 0644)
	if err != nil {
		return fmt.Errorf("error writing YAML file: %v", err)
	}

	fmt.Printf("\nEnvironment configuration saved to %s\n", filename)
	return nil
}

var captureCmd = &cobra.Command{
	Use:   "capture",
	Short: "Capture the current environment.",
	Long:  `Captures the current relevant environment state and stores it in a configuration YAML file.`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			return
		}

		// Capture different components of the environment
		dotfiles, err := captureDotFiles(dir)
		if err != nil {
			fmt.Printf("Error capturing dot files: %v\n", err)
			return
		}

		envVars, err := captureEnvVariables()
		if err != nil {
			fmt.Printf("Error capturing environment variables: %v\n", err)
			return
		}

		// Build and save schema
		schema := buildEnvironmentSchema(dotfiles, envVars)
		if err := saveEnvironmentSchema(schema, "environment.yaml"); err != nil {
			fmt.Printf("Error saving environment configuration: %v\n", err)
			return
		}
	},
}
