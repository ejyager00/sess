/*
Copyright © 2024 Eric Yager
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ejyager00/sess/internal/io"
	"github.com/ejyager00/sess/internal/models"
	"github.com/spf13/cobra"
)

type extensionInstaller interface {
	installExtension(extension models.Extension) error
	isAvailable() bool
	getName() string
}

type vsCodeInstaller struct {
	command string
}

func (v vsCodeInstaller) installExtension(extension models.Extension) error {
	cmd := exec.Command(v.command, "--install-extension", extension.Id)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install extension %s: %v\n%s", extension.Id, err, output)
	}
	fmt.Printf("Successfully installed extension %s\n", extension.Id)
	return nil
}

func (v vsCodeInstaller) isAvailable() bool {
	_, err := exec.LookPath(v.command)
	return err == nil
}

func (v vsCodeInstaller) getName() string {
	return strings.TrimPrefix(v.command, "--")
}

func installExtensions(extensions []models.Extension) error {
	// Map of supported IDE installers
	installers := []extensionInstaller{
		vsCodeInstaller{command: "code"},
		vsCodeInstaller{command: "codium"},
	}

	// Group extensions by IDE
	extensionsByIde := make(map[string][]models.Extension)
	for _, extension := range extensions {
		extensionsByIde[extension.Ide] = append(extensionsByIde[extension.Ide], extension)
	}

	// Install extensions for each IDE
	for ide, exts := range extensionsByIde {
		var installer extensionInstaller
		found := false

		// Find first available installer for this IDE
		for _, i := range installers {
			if strings.EqualFold(ide, "vscode") && i.isAvailable() {
				installer = i
				found = true
				break
			}
		}

		if !found {
			fmt.Printf("Skipping extensions for %s - no compatible installer found\n", ide)
			continue
		}

		fmt.Printf("Installing extensions using %s...\n", installer.getName())
		for _, extension := range exts {
			if err := installer.installExtension(extension); err != nil {
				return err
			}
		}
	}
	return nil
}

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

		if err := installExtensions(schema.Extensions); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
