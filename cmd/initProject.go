package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ubavic/mint/schema"
	"gopkg.in/yaml.v3"
)

type InitProjectConfig struct {
	name   string
	author string
}

func initProject() error {
	config, err := initProjectParseConfig()
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
		Version: "v0.0.1",
		Name:    config.name,
		Author:  config.author,
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

	printSuccess("Project initialized successfully")
	fmt.Println("You can build your project with:")
	fmt.Println("mint build")

	return nil
}

func initProjectParseConfig() (InitProjectConfig, error) {
	initProjectFlags := InitProjectConfig{}

	initFlagSet := flag.NewFlagSet("init", flag.ExitOnError)
	nameFlag := initFlagSet.String("name", "", "Project name")
	authorFlag := initFlagSet.String("author", "", "Project author")

	initFlagSet.Usage = printInitHelp

	err := initFlagSet.Parse(os.Args[2:])
	if err != nil {
		return initProjectFlags, fmt.Errorf("parsing flags: %w", err)
	}

	initProjectFlags.name = *nameFlag
	initProjectFlags.author = *authorFlag

	if initProjectFlags.name == "" {
		wd, _ := os.Getwd()
		initProjectFlags.name = filepath.Base(wd)
	}

	if initProjectFlags.author == "" {
		initProjectFlags.author = os.Getenv("USER")
	}

	return initProjectFlags, nil
}

func printInitHelp() {
	fmt.Println("Usage: mint init [options]")
	fmt.Println("Options:")
	fmt.Println("  -name NAME - Project name. Optional. If not provided, the current directory name is used.")
	fmt.Println("  -author AUTHOR - Project author. Optional. If not provided, the current user is used.")
	fmt.Println("  -h, --help - Print this help message")
}
