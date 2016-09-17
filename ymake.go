package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/daviddengcn/go-colortext"

	. "github.com/nananas/ymake/ymakefile"
)

// use go linker tool '-ldflags "-X main.VERSION=..."' when building to set this value
var VERSION string

func main() {

	var f_version bool
	var f_debug bool
	var f_verbose bool
	var f_force bool
	var f_printonly bool

	flag.BoolVar(&f_version, "V", false, "shorthand for -version")
	flag.BoolVar(&f_version, "version", false, "Prints version number.")

	flag.BoolVar(&f_debug, "d", false, "shorthand for -debug")
	flag.BoolVar(&f_debug, "debug", false, "Toggles the debug flag on.")

	flag.BoolVar(&f_verbose, "v", false, "shorthand for -verbose")
	flag.BoolVar(&f_verbose, "verbose", false, "Toggles the verbose flag on.")

	flag.BoolVar(&f_force, "f", false, "shorthand for -force")
	flag.BoolVar(&f_force, "force", false, "Execute commands, even if dependencies are met.")

	flag.BoolVar(&f_printonly, "p", false, "shorthand for -printonly")
	flag.BoolVar(&f_printonly, "printonly", false, "Only print the commands, do not execute them. However, variable commands are still executed.")

	flag.Usage = printHelp

	flag.Parse()

	if f_version {
		fmt.Println("Ymake version v" + VERSION)
		return
	}

	DEBUG = f_debug
	VERBOSE = f_verbose
	FORCE = f_force
	PRINTONLY = f_printonly

	f, err := os.Open("./Ymakefile")
	if err != nil {
		log.Fatal(err)
	}

	Print("[ ] Start")

	var YMAKE *YMakefile
	var VARIABLES *Variables
	var SHELL string

	YMAKE, VARIABLES, SHELL = LoadConfig(f)

	containsBlockName := false

	for _, a := range flag.Args()[0:] {
		if a[0] != '-' {
			containsBlockName = true
			exists, err := RunBlock(flag.Args()[0], YMAKE, VARIABLES, SHELL, 0)
			if !exists {
				Error("No block found called " + flag.Args()[0])
			}
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if !containsBlockName {
		_, err := RunBlock("default", YMAKE, VARIABLES, SHELL, 0)
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
	Ymake is a build tool with the same purpose as GNU make:
		- making the building of software projects easier.
		- combining useful commands into a single file (clean, test, package...)
		- use dependencies to remove the need to recompile from scratch
`)
	ct.Foreground(ct.White, true)
	fmt.Print(`	Usage:`)
	ct.ResetColor()
	fmt.Println(`
		ymake [options] [blockname]
`)
	ct.Foreground(ct.White, true)
	fmt.Print(`	The options:`)
	ct.ResetColor()
	fmt.Println(`
		ymake [options] [blockname]
		-h, -help   		Prints this help text
		-V, -version		Prints the version
		-v, -verbose		The output will be more verbose
		-d, -debug  		The debug output contains more info
		-p, -printonly		Only print commands, do not run.
		-f, -force  		Force execution of commands, even if deps are met.

		[*] Every option can use either - or --

`)
	ct.Foreground(ct.White, true)
	fmt.Print(`	The blockname:`)
	ct.ResetColor()
	fmt.Println(`
		is the name of any block specified under 'targets'.
		When ommitted, the "default" block will be executed.
`)

}
