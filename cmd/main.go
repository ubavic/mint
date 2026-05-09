package main

import (
	_ "embed"
	"fmt"
	"os"
)

//go:embed version
var version string

func main() {
	if len(os.Args) < 2 {
		printError("no command provided")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		err := initProject()
		if err != nil {
			printError(err.Error())
			os.Exit(2)
		}
	case "build":
		err := buildProject()
		if err != nil {
			printError(err.Error())
			os.Exit(3)
		}
	case "version":
		fmt.Println(version)
	case "help":
		printHelp()
	default:
		printError(fmt.Sprintf("unknown command: %s", os.Args[1]))
		os.Exit(4)
	}
}

func printHelp() {
	command := ""
	if len(os.Args) > 2 {
		command = os.Args[2]
	}

	switch command {
	case "init":
		printInitHelp()
		return
	case "build":
		printBuildHelp()
		return
	}

	fmt.Println("Usage: mint <command> [options]")
	fmt.Println("Commands:")
	fmt.Println("  init - Initialize a new Mint project")
	fmt.Println("  build - Build a Mint project")
	fmt.Println("  version - Print the Mint version")
	fmt.Println("  help - Print this help message")
}
