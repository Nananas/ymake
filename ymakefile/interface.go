package ymakefile

import "fmt"

var VERBOSE = false
var DEBUG = false
var FORCE = false

func Print(in string) {
	fmt.Println("\t" + in)
}

func PrintV(in string) {
	if VERBOSE {
		fmt.Println("\t\t" + in)
	}
}
