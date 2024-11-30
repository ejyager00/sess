/*
Copyright © 2024 Eric Yager
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
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

// buildEnvironmentSchema constructs the environment schema from captured data
func buildEnvironmentSchema(dotfiles []string) *models.EnvironmentSchema {
	return &models.EnvironmentSchema{
		Dotfiles: dotfiles,
		// Future fields will be added here:
		// Tools: capturedTools,
		// EnvVariables: capturedEnvVars,
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

		// Build and save schema
		schema := buildEnvironmentSchema(dotfiles)
		if err := saveEnvironmentSchema(schema, "environment.yaml"); err != nil {
			fmt.Printf("Error saving environment configuration: %v\n", err)
			return
		}
	},
}
