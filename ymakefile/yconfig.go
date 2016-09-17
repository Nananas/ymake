package ymakefile

import (
	"io"
	"io/ioutil"
	"log"
	"reflect"

	yaml "gopkg.in/yaml.v2"
)

type YMakefile struct {
	Blocks struct {
		Default *YBlock            ",omitempty"
		Others  map[string]*YBlock ",inline"
	}

	Options struct {
		// Cores int ",omitempty"
		Shell string ",omitempty"
	}

	Variables struct {
		V map[string]interface{} ",omitempty,inline"
	}
}

type YBlock struct {
	Cmd     Either
	Post    Either ",omitempty"
	Target  string ",omitempty"
	Deps    Either ",omitempty"
	Pattern string ",omitempty"
	Stdin   string ",omitempty"

	// options
	Parallel bool ",omitempty"
	Hide     bool ",omitempty"
}

// Either blocks are either:
// 	- a single string
// 	- a list of strings
//
type Either interface{}
type EitherAction func(s string) bool

type Variables map[string]string

var (
	DEFAULT_VARIABLES = []string{
		"WD",
		"HOME",
	}

	DEFAULT_VARIABLE_COMMANDS = []string{
		"pwd",
		"echo $HOME",
	}
)

// Utility function to run a function over every string in an Either block
// Returns false if something happened
//
func HandleEither(either Either, action EitherAction) bool {

	switch t := either.(type) {
	case string:
		return action(t)
	case []interface{}: //[]string
		for _, e := range t {

			if e == nil || reflect.TypeOf(e).Kind() != reflect.String {
				continue
			}

			s := e.(string)

			if !action(s) {
				return false
			}

		}

		return true
	default:
		log.Println("This case should not happen")
		return false

	}

}

// Load ymakefile from reader
//
func LoadConfig(reader io.Reader) (*YMakefile, *Variables, string) {

	if DEBUG {
		log.SetFlags(log.Lshortfile)
	} else {
		log.SetFlags(0)
	}

	var config YMakefile

	// read ymakefile
	//
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}

	// parse ymakefile
	//
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		Error("[E] makefile error!")
		log.Fatal(err)
	}

	// parse options
	//
	shell := config.Options.Shell

	// parse variables
	//
	variables := make(Variables, len(config.Variables.V)+len(DEFAULT_VARIABLES))
	for i, e := range DEFAULT_VARIABLES {
		output, err := ExecuteCapture(DEFAULT_VARIABLE_COMMANDS[i], shell)
		if err != nil {
			log.Fatal(output, err)
		}

		variables[e] = output
		PrintV("[d] " + e + ": " + output)
	}

	for k, v := range config.Variables.V {
		switch cmd := v.(type) {
		case string:

			output, err := ExecuteCapture(cmd, shell)
			if err != nil {
				log.Fatal(output, err)
			}

			variables[k] = output

			// fmt.Print("\t)
			PrintVariable(k + ": " + output)
		default:
			log.Fatal("Variable " + k + " is not a string")
		}
	}

	return &config, &variables, shell
}
