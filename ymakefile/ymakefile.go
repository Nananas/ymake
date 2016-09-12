package ymakefile

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Pattern struct {
	// contains the group results of the regexp:
	// $1,$2 ...
	List *[]string
}

// Replaces every known variable with its value
//
func Vars(in string, variables *Variables) string {
	for k, v := range *variables {
		in = strings.Replace(in, "{"+k+"}", v, -1)
	}

	return in
}

// Replaces every known pattern with its value
// $0 is the list of every pattern joined by a space
func Patterns(in string, patterns *[]string) string {

	if patterns == nil {
		return in
	}

	out := in

	for i, p := range *patterns {
		out = strings.Replace(out, "$"+strconv.Itoa(i+1), p, -1)
	}

	allpatterns := strings.Join(*patterns, " ")

	out = strings.Replace(out, "$0", allpatterns, -1)

	return out
}

// Converts a simple Glob pattern to a Regexp
//
func createRegex(s string) *regexp.Regexp {
	expr := s

	expr = strings.Replace(expr, "*", "[\\w\\W]*", -1)
	expr = strings.Replace(expr, ".", "\\.", -1)
	expr = "^" + expr + "$"

	// expr = strings.Replace(expr, "~", home_dir, -1)

	regex, err := regexp.Compile(expr)
	if err != nil {
		log.Fatal(err)
	}

	return regex
}

// Runs a block with blockname
// returns false if the block does not exist
// returns error if something happend before execution
// execution errors through stderr are not returned
// intent handles the indention of this block as subblock to a parent
//
func RunBlock(blockname string, ymakefile *YMakefile, variables *Variables, shell string, indent int) (bool, error) {

	var block *YBlock

	if blockname != "default" {
		if _, ok := ymakefile.Blocks.Others[blockname]; !ok {
			return false, nil
		}

		block = ymakefile.Blocks.Others[blockname]
	} else {
		block = ymakefile.Blocks.Default

	}

	if block == nil {
		Error("[E] No block found with name '" + blockname + "'")
		return false, nil
	}

	PrintBlock(blockname, indent)

	patterns := make([]Pattern, 0)

	if block.Pattern != "" {

		patt := Vars(block.Pattern, variables)

		clean_line := strings.Replace(patt, "(", "", -1)
		clean_line = strings.Replace(clean_line, ")", "", -1)

		regex := createRegex(patt)

		// find all files matching the clean pattern
		matches, err := filepath.Glob(clean_line)
		if err != nil {
			log.Println(err)
		}

		// for each file, match with the original regexp to find group results
		for _, patterned_files := range matches {

			regexmatch := regex.FindAllStringSubmatch(patterned_files, -1)

			// Get submatch groups from regex of the matching file line
			//
			list := []string{}
			for _, m := range regexmatch {
				list = append(list, m[1])
			}

			patterns = append(patterns, Pattern{List: &list})
		}
	}

	// no pattern means no capture, no $-extention
	// but the patterns loop has to run at least once though
	if len(patterns) == 0 {
		patterns = append(patterns, Pattern{List: nil})
	}

	for _, p := range patterns {
		shouldRun := false

		if block.Deps != nil {
			if block.Target == "" {
				return true, errors.New("No target specified, but found dependencies. Maybe use 'pattern'?")
			}

			deps_files := make([]string, 0)

			if !HandleEither(block.Deps, func(dep string) bool {

				dep = Vars(Patterns(dep, p.List), variables)

				matches, err := filepath.Glob(dep)
				if err != nil {
					log.Println(err)
					return false
				}

				deps_files = append(deps_files, matches...)

				return true
			}) {
				log.Println("Exit deps")
				return true, errors.New("Dependency error")
			}

			(*variables)["DEPS"] = strings.Join(deps_files, " ")

			target := Vars(Patterns(block.Target, p.List), variables)
			stat, err := os.Stat(target)
			if err != nil {
				Print("\t[?] " + target)
				shouldRun = true
			} else {
				for _, d := range deps_files {

					dst, err := os.Stat(d)
					if err != nil {
						log.Fatal(err)
					}

					if dst.ModTime().Sub(stat.ModTime()).Seconds() > 0 {
						shouldRun = true
						Print("[M] " + dst.Name())
					} else {
						PrintV("[s] " + dst.Name())
					}
				}
			}

		} else {
			shouldRun = true
		}

		if shouldRun || FORCE {

			if block.Parallel {
				// PARALLEL
				//

				errorchan := make(chan error, 1)
				if block.Cmd == nil {

				} else {
					switch t := block.Cmd.(type) {
					case string:
						return true, errors.New("Stop. Expected multiple commands when running parallel!")
					case []interface{}:
						// channel of at least the size of all possible commands
						//
						donechan := make(chan bool, len(t))

						// keeps track of how many commands there actually are
						// ignoring empty strings etc.
						//
						cmdcount := 0

						for _, e := range t {

							if e == nil || reflect.TypeOf(e).Kind() != reflect.String {
								continue
							}

							cmdcount++

							s := e.(string)

							cmd := Vars(Patterns(s, p.List), variables)
							if !block.Hide {
								PrintCmd(cmd, indent)
							}

							stdin := block.Stdin

							// execute command in parallel
							//
							ExecuteStdParallel(cmd, stdin, shell, errorchan, donechan)
						}

						// check if all child processes are completed
						//
						for i := 0; i < cmdcount; i++ {
							<-donechan
						}

						// check if an error occured
						//
						select {
						case err := <-errorchan:
							ErrorCmd(err.Error(), indent)
							return true, err
						default:
							// nothing bad happend
						}
					default:
						return false, errors.New("This case should not happen")
					}
				}

			} else {
				// SEQUENTIAL
				//
				if block.Cmd == nil {

				} else {
					if !HandleEither(block.Cmd, func(cmd string) bool {
						cmd = Vars(Patterns(cmd, p.List), variables)
						if !block.Hide {
							PrintCmd(cmd, indent)
						}

						stdin := block.Stdin

						err := ExecuteStd(cmd, stdin, shell)
						if err != nil {
							ErrorCmd(err.Error(), indent)
							return false
						}

						return true
					}) {
						return true, errors.New("Stop.")
					}
				}

			}

		}

	}

	if block.Post != nil {
		HandleEither(block.Post, func(s string) bool {
			exists, err := RunBlock(s, ymakefile, variables, shell, indent+1)
			if !exists {
				Error("No block found called " + s)
				return false
			}
			if err != nil {
				PrintInfo(err.Error(), indent)
				return false
			}

			return true
		})
	}

	return true, nil
}
