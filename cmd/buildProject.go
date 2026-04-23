package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ubavic/mint/schema"
)

type BuildProjectFlags struct {
	Target string
	Schema string
	JSON   bool
}

func buildProject() error {
	flags, err := buildProjectParseFlags()
	if err != nil {
		return err
	}

	_, err = schema.OpenAndValidateSchema(flags.Schema, version)
	if err != nil {
		return err
	}

	return nil
}

func buildProjectParseFlags() (BuildProjectFlags, error) {
	buildProjectFlags := BuildProjectFlags{}

	buildFlagSet := flag.NewFlagSet("build", flag.ExitOnError)
	target := buildFlagSet.String("target", "", "Select target from schema. If not provided a first target from schema is used.")
	schemaFileFlag := buildFlagSet.String("schema", "", "Specifies a schema file.")
	jsonFlag := buildFlagSet.Bool("json", false, "Output JSON AST instead target format")

	err := buildFlagSet.Parse(os.Args[2:])
	if err != nil {
		return buildProjectFlags, fmt.Errorf("parsing flags: %w", err)
	}

	buildProjectFlags.Target = *target
	buildProjectFlags.Schema = *schemaFileFlag
	buildProjectFlags.JSON = *jsonFlag

	return buildProjectFlags, nil
}
