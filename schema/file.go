package schema

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const DefaultSchemaFilename = "mint.yaml"

func OpenAndValidateSchema(filename string) (Schema, error) {
	var schema Schema

	if filename == "" {
		filename = DefaultSchemaFilename
	}

	schemaFile, err := os.ReadFile(filename)
	if err != nil {
		return schema, fmt.Errorf("can't open schema file \"%s\": %w", filename, err)
	}

	err = yaml.Unmarshal(schemaFile, &schema)
	if err != nil {
		return schema, fmt.Errorf("can't unmarshal schema: %w", err)
	}

	err = schema.Check()
	if err != nil {
		return schema, fmt.Errorf("validating schema: %w", err)
	}

	return schema, nil
}
