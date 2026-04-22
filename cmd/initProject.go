package main

import (
	"fmt"
	"os"

	"github.com/ubavic/mint/schema"
	"gopkg.in/yaml.v3"
)

type InitProjectFlags struct{}

func initProject() error {
	_, err := initProjectParseFlags()
	if err != nil {
		return err
	}

	_, err = os.Stat(schema.DefaultSchemaFilename)
	if err == nil {
		return fmt.Errorf("schema file already exists in current folder")
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("file error %w", err)
	}

	s := schema.Schema{
		Mint:    version,
		Version: "v0.0.0",
		Author:  os.Getenv("USER"),
		Source: schema.Source{
			AllowedRootCommands: "p",
			Commands: []schema.Command{
				{
					Command:     "p",
					Description: "paragraph",
				},
			},
		},
		Targets: []schema.Target{
			{
				Name:      "Text",
				Extension: "txt",
			},
		},
	}

	yml, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Errorf("marshaling schema: %w", err)
	}

	err = os.WriteFile(schema.DefaultSchemaFilename, yml, 0600)
	if err != nil {
		return fmt.Errorf("writing schema: %w", err)
	}

	return nil
}

func initProjectParseFlags() (InitProjectFlags, error) {
	initProjectFlags := InitProjectFlags{}

	return initProjectFlags, nil
}
