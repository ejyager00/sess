/*
Copyright © 2024 Eric Yager
*/
package models

type EnvironmentSchema struct {
	Tools        []Tool            `yaml:"tools"`
	EnvVariables map[string]string `yaml:"env_variables"`
	Dotfiles     []string          `yaml:"dotfiles"`
	Extensions   []Extension       `yaml:"extensions"`
}

type Tool struct {
	Name                string `yaml:"name"`
	Version             string `yaml:"version"`
	InstallCommand      string `yaml:"install_command"`
	VersionCheckCommand string `yaml:"version_check_command"`
}

type Extension struct {
	Id      string `yaml:"id"`
	Ide     string `yaml:"ide"`
	Version string `yaml:"version"`
}
