package main

import (
	"fmt"
	"log"
	"os"

	. "github.com/nananas/ymake/ymakefile"
)

// use go linker tool '-ldflags "-X main.VERSION=..."' when building to set this value
var VERSION string

func main() {

	if contains(os.Args, "-V", "--version") {
		fmt.Println("Ymake version v" + VERSION)
		return
	}

	if contains(os.Args, "-h", "--help") {
		printHelp()
		return
	}

	if contains(os.Args, "-d", "--debug") {
		DEBUG = true
	}

	if contains(os.Args, "-v", "--verbose") {
		VERBOSE = true
	}

	if contains(os.Args, "-f", "--force") {
		FORCE = true
	}

	f, err := os.Open("./Ymakefile")
	if err != nil {
		log.Fatal(err)
	}

	Print("[ ] Start")

	var YMAKE *YMakefile
	var VARIABLES *Variables
	YMAKE, VARIABLES = LoadConfig(f)

	containsBlockName := false

	for _, a := range os.Args[1:] {
		if a[0] != '-' {
			containsBlockName = true
			exists, err := RunBlock(os.Args[1], YMAKE, VARIABLES)
			if !exists {
				fmt.Println("No block found called " + os.Args[1])
			}
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if !containsBlockName {
		_, err := RunBlock("default", YMAKE, VARIABLES)
		if err != nil {
			log.Fatal(err)
		}
	}

	Print("[ ] Exit")
}

func contains(list []string, targets ...string) bool {
	for _, e := range list {
		for _, t := range targets {
			if e == t {
				return true
			}
		}
	}

	return false
}

func notAny(in string, any ...string) bool {
	for _, a := range any {
		if in == a {
			return false
		}
	}

	return true
}

func printHelp() {
	fmt.Println(`
Ymake is a build tool analog to GNU make.

Usage:
	ymake [options|blockname]	

The options:
	-h --help 			Prints this help text
	-V --version 		Prints the version
	-v --verbose 		The output will be more verbose 
	-d --debug 			The debug output contains more info

The blockname:
	is the name of any block specified under 'targets'. 
	When ommitted, the default will be executed.
`)
}
