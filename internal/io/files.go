/*
Copyright © 2024 Eric Yager
*/
package io

import (
	"fmt"
	"os"

	"github.com/ejyager00/sess/internal/models"
	"github.com/goccy/go-yaml"
)

func ParseYamlFile(yamlFilePath string) (*models.EnvironmentSchema, error) {
	data, err := os.ReadFile(yamlFilePath)
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
