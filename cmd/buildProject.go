package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ubavic/mint/parser"
	"github.com/ubavic/mint/schema"
	"github.com/ubavic/mint/writer"
)

type BuildProjectConfig struct {
	Target    string
	Schema    string
	JSON      bool
	InputFile string
}

func buildProject() error {
	flags, err := buildProjectParseFlags()
	if err != nil {
		return err
	}

	currentSchema, err := schema.OpenAndValidateSchema(flags.Schema, version)
	if err != nil {
		return err
	}

	var target *schema.Target

	if flags.Target != "" {
		target, err = currentSchema.GetTarget(flags.Target)
		if err != nil {
			return fmt.Errorf("can't find target: %w", err)
		}
	} else {
		if len(currentSchema.Targets) == 0 {
			return fmt.Errorf("no targets found in schema")
		}
		target = &currentSchema.Targets[0]
	}

	outDir := filepath.Join("output", target.Name)
	err = os.MkdirAll(outDir, 0755)
	if err != nil {
		return fmt.Errorf("can't create output directory: %w", err)
	}

	if flags.InputFile != "" {
		printInfo("Processing file: " + flags.InputFile)

		err = processFile(flags.InputFile, outDir, &currentSchema, target, flags.JSON)
		if err != nil {
			return err
		}

		printSuccess(fmt.Sprintf("File %q processed successfully", flags.InputFile))

		return nil
	}

	fmt.Println("here")

	root := "."
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".atex" {
			return nil
		}

		err = processFile(path, outDir, &currentSchema, target, flags.JSON)
		if err != nil {
			return fmt.Errorf("error while processing file %q: %w", path, err)
		}

		return nil

	})
	if err != nil {
		return err
	}

	printSuccess("All .atex files processed successfully")

	return nil
}

func processFile(path string, outDir string, schema *schema.Schema, target *schema.Target, json bool) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("can't open file %q: %w", path, err)
	}

	fileBuf := bufio.NewReader(file)
	tokenizer := parser.NewTokenizer(fileBuf)

	tokens := tokenizer.Tokenize()

	schemaParser := parser.NewParser(tokens, schema)
	doc, err := schemaParser.Parse()
	if err != nil {
		return fmt.Errorf("error while parsing %q: %w", path, err)
	}

	var rendered string
	if json {
		rendered = string(doc.Json())
	} else {
		rendered = writer.Write(target, doc)
	}

	err = os.WriteFile(filepath.Join(outDir, filepath.Base(path)), []byte(rendered), 0644)
	if err != nil {
		return fmt.Errorf("can't write to file %q: %w", filepath.Join(outDir, filepath.Base(path)), err)
	}

	return nil
}

func buildProjectParseFlags() (BuildProjectConfig, error) {
	buildProjectFlags := BuildProjectConfig{}

	buildFlagSet := flag.NewFlagSet("build", flag.ExitOnError)
	target := buildFlagSet.String("target", "", "Select target from schema. If not provided a first target from schema is used.")
	schemaFileFlag := buildFlagSet.String("schema", "", "Specifies a schema file.")
	jsonFlag := buildFlagSet.Bool("json", false, "Output JSON AST instead target format")
	inputFileFlag := buildFlagSet.String("input", "", "Specifies a input file.")

	buildFlagSet.Usage = printBuildHelp

	err := buildFlagSet.Parse(os.Args[2:])
	if err != nil {
		return buildProjectFlags, fmt.Errorf("parsing flags: %w", err)
	}

	buildProjectFlags.Target = *target
	buildProjectFlags.Schema = *schemaFileFlag
	buildProjectFlags.JSON = *jsonFlag
	buildProjectFlags.InputFile = *inputFileFlag

	return buildProjectFlags, nil
}

func printBuildHelp() {
	fmt.Println("Usage: mint build [options]")
	fmt.Println("Options:")
	fmt.Println("  -target TARGET - Select target from schema. If not provided a first target from schema is used.")
	fmt.Println("  -schema SCHEMA - Specifies a schema file. If not provided a schema in current directory is used.")
	fmt.Println("  -json - Output JSON AST instead target format")
	fmt.Println("  -input INPUT - Specifies a input file. If not provided all .atex files in current directory are processed.")
	fmt.Println("  -h, --help - Print this help message")
}
