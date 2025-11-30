package main

import (
	"bufio"
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/ubavic/mint/parser"
	"github.com/ubavic/mint/schema"
	"github.com/ubavic/mint/writer"
	"gopkg.in/yaml.v3"
)

//go:embed version
var version string

func main() {
	inputFileFlag := flag.String("in", "", "Specifies a input file")
	schemaFileFlag := flag.String("schema", "", "Specifies a schema file")
	targetFlag := flag.String("target", "", "Select target from schema")
	outputFileFlag := flag.String("out", "", "Specifies a output file")
	jsonFlag := flag.Bool("json", false, "Output JSON")
	versionFlag := flag.Bool("version", false, "Print version and exit")

	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		return
	}

	config := Config{
		InputFile:  *inputFileFlag,
		SchemaFile: *schemaFileFlag,
		Target:     *targetFlag,
		Json:       *jsonFlag,
		OutputFile: *outputFileFlag,
	}

	err := process(config)
	if err != nil {
		printErrorAndExit(err)
	}
}

type Config struct {
	InputFile  string
	SchemaFile string
	Target     string
	Json       bool
	OutputFile string
}

func process(config Config) error {
	if config.InputFile == "" {
		return fmt.Errorf("expected input file")
	}

	file, err := os.Open(config.InputFile)
	if err != nil {
		return fmt.Errorf("can't open file \"%s\": %v", config.InputFile, err.Error())
	}

	var newSchema *schema.Schema = nil

	if config.SchemaFile != "" {
		schemaFile, err := os.ReadFile(config.SchemaFile)
		if err != nil {
			return fmt.Errorf("can't open file \"%s\": %v", config.SchemaFile, err.Error())
		}

		err = yaml.Unmarshal(schemaFile, &newSchema)
		if err != nil {
			return fmt.Errorf("can't unmarshal schema: %v", err.Error())
		}

		err = newSchema.Check()
		if err != nil {
			return fmt.Errorf("schema validation: %w", err)
		}
	}

	fileBuf := bufio.NewReader(file)
	tokenizer := parser.NewTokenizer(fileBuf)

	tokens := tokenizer.Tokenize()

	parser := parser.NewParser(tokens, newSchema)
	doc, err := parser.Parse()
	if err != nil {
		return fmt.Errorf("error while parsing \"%s\": %v", config.InputFile, err.Error())
	}

	var target *schema.Target = nil
	if config.Target != "" {

		if newSchema == nil {
			return fmt.Errorf("schema is not loaded")
		}

		target, err = newSchema.GetTarget(config.Target)
		if err != nil {
			return fmt.Errorf("can't find target: %v", err.Error())
		}
	}

	rendered := ""

	if config.Json {
		rendered = string(doc.Json())
	} else {
		rendered = writer.Write(target, doc)
	}

	if config.OutputFile != "" {
		err = os.WriteFile(config.OutputFile, []byte(rendered), 0644)
		if err != nil {
			return fmt.Errorf("can't write to file \"%s\": %v", config.OutputFile, err.Error())
		}
	} else {
		fmt.Println(rendered)
	}

	return nil
}

func printErrorAndExit(err error) {
	fmt.Println("\x1b[91m" + err.Error() + "\x1b[0m")
	os.Exit(1)
}
