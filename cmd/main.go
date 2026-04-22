package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
)

//go:embed version
var version string

func main() {
	if len(os.Args) < 2 {
		printErrorAndExit(fmt.Errorf("no command or flag provided"))
	}

	switch os.Args[1] {
	case "init":
		err := initProject()
		if err != nil {
			printErrorAndExit(err)
		}
		return
	case "build":
		err := buildProject()
		if err != nil {
			printErrorAndExit(err)
		}
		return
	}

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

func printErrorAndExit(err error) {
	fmt.Println("\x1b[91m" + err.Error() + "\x1b[0m")
	os.Exit(1)
}
