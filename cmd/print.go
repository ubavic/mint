package main

import "fmt"

func printSuccess(message string) {
	fmt.Println("\x1b[92m" + message + "\x1b[0m")
}

func printError(message string) {
	fmt.Println("\x1b[91m" + message + "\x1b[0m")
}

func printWarning(message string) {
	fmt.Println("\x1b[93m" + message + "\x1b[0m")
}

func printInfo(message string) {
	fmt.Println("\x1b[94m" + message + "\x1b[0m")
}
