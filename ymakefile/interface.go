package ymakefile

import (
	"github.com/daviddengcn/go-colortext"

	"fmt"
)

var VERBOSE = false
var DEBUG = false
var FORCE = false
var PRINTONLY = false

func Print(in string) {
	fmt.Println("\t" + in)
}

func PrintV(in string) {
	if VERBOSE {
		fmt.Println("\t\t" + in)
	}
}

func PrintInfo(in string, indent int) {
	fmt.Print("\t")
	for i := 0; i < indent; i++ {
		fmt.Print("   ")
	}
	ct.Foreground(ct.Blue, false)
	fmt.Print("[-] ")
	ct.ResetColor()
	fmt.Println(in)

}

func PrintCmd(in string, indent int) {
	fmt.Print("\t")
	for i := 0; i < indent; i++ {
		fmt.Print("   ")
	}
	ct.Foreground(ct.Green, false)
	fmt.Print("[>] ")
	ct.ResetColor()
	fmt.Println(in)
}

func ErrorCmd(in string, indent int) {
	fmt.Print("\t")
	for i := 0; i < indent; i++ {
		fmt.Print("   ")
	}
	ct.Foreground(ct.Red, true)
	fmt.Print("[<] ")
	ct.ResetColor()
	fmt.Println(in)
}

func PrintBlock(in string, indent int) {
	fmt.Print("\t")
	for i := 0; i < indent; i++ {
		fmt.Print("   ")
	}
	ct.Foreground(ct.Cyan, false)
	fmt.Print("[B] ")
	ct.ResetColor()
	fmt.Println(in)
}

func PrintVariable(in string) {
	ct.Foreground(ct.Yellow, false)
	fmt.Print("\t[V] ")
	ct.ResetColor()
	fmt.Println(in)
}

func Error(in string) {
	ct.Foreground(ct.Red, true)
	fmt.Println("\t" + in)
	ct.ResetColor()
}
