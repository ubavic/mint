package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/ubavic/mint/parser"
	"github.com/ubavic/mint/schema"
	"github.com/ubavic/mint/writer"
	"gopkg.in/yaml.v3"
)

func main() {
	inputFileFlag := flag.String("in", "", "Specifies a input file")
	schemaFileFlag := flag.String("schema", "", "Specifies a schema file")
	targetFlag := flag.String("target", "", "Select target from schema")
	flag.Parse()

	if *inputFileFlag == "" {
		fmt.Println("Expected input file")
		return
	}

	if *schemaFileFlag == "" {
		fmt.Println("Expected schema file")
		return
	}

	schemaFile, err := os.ReadFile(*schemaFileFlag)
	if err != nil {
		fmt.Printf("Can't open file \"%s\": %v", *schemaFileFlag, err.Error())
		return
	}

	file, err := os.Open(*inputFileFlag)
	if err != nil {
		fmt.Printf("Can't open file \"%s\": %v", *inputFileFlag, err.Error())
		return
	}

	var newSchema schema.Schema

	err = yaml.Unmarshal(schemaFile, &newSchema)
	if err != nil {
		fmt.Printf("Can't unmarshal schema: %v", err.Error())
		return
	}

	fileBuf := bufio.NewReader(file)
	tokenizer := parser.NewTokenizer(fileBuf)

	tokens := tokenizer.Tokenize()

	parser := parser.NewParser(tokens, newSchema)
	doc, err := parser.Parse()
	if err != nil {
		fmt.Printf("Error while parsing \"%s\": %v", *inputFileFlag, err.Error())
		return
	}

	target, err := newSchema.GetTarget(*targetFlag)
	if err != nil {
		fmt.Printf("Can't find target: %v", err.Error())
		return
	}

	rendered := writer.Write(target, doc)

	fmt.Println(rendered)
}
