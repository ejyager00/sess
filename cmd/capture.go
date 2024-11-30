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

// getDotFilesSelection combines listing and selection of dot files
func getDotFilesSelection(dir string) ([]string, error) {
	dotFiles, err := listDotFiles(dir)
	if err != nil {
		return nil, err
	}

	selectedFiles, err := promptForDotFileSelection(dotFiles)
	if err != nil {
		return nil, err
	}

	return selectedFiles, nil
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

		selectedFiles, err := getDotFilesSelection(dir)
		if err != nil {
			fmt.Printf("Error selecting dot files: %v\n", err)
			return
		}

		fmt.Println("\nSelected dot files:")
		for _, file := range selectedFiles {
			fmt.Println(file)
		}

		// Create environment schema with selected dotfiles
		schema := models.EnvironmentSchema{
			Dotfiles: selectedFiles,
		}

		// Marshal to YAML
		yamlData, err := yaml.Marshal(&schema)
		if err != nil {
			fmt.Printf("Error marshaling YAML: %v\n", err)
			return
		}

		// Write to environment.yaml file
		err = os.WriteFile("environment.yaml", yamlData, 0644)
		if err != nil {
			fmt.Printf("Error writing YAML file: %v\n", err)
			return
		}

		fmt.Println("\nEnvironment configuration saved to environment.yaml")
	},
}
