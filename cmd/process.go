package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ubavic/mint/parser"
	"github.com/ubavic/mint/schema"
	"github.com/ubavic/mint/writer"
)

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
		s, err := schema.OpenAndValidateSchema(config.SchemaFile, version)
		newSchema = &s
		if err != nil {
			return err
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
