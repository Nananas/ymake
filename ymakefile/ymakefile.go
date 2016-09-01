package ymakefile

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
//
func RunBlock(blockname string, ymakefile *YMakefile, variables *Variables) (bool, error) {

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
		Print("[E] No block found with name '" + blockname + "'")
		return false, nil
	}

	Print("[B] " + blockname)

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

			if block.Cmd == nil {

			} else {
				if !HandleEither(block.Cmd, func(cmd string) bool {
					cmd = Vars(Patterns(cmd, p.List), variables)
					if !block.Hide {
						Print("[>] " + cmd)
					}
					err := ExecuteStd(cmd)
					if err != nil {
						Print("[<] " + err.Error())
						return false
					}

					return true
				}) {
					return true, errors.New("Stop.")
				}
			}

		}

	}

	if block.Post != nil {
		HandleEither(block.Post, func(s string) bool {
			exists, err := RunBlock(s, ymakefile, variables)
			if !exists {
				fmt.Println("No block found called " + s)
				return false
			}
			if err != nil {
				log.Println("\t--- " + err.Error())
				return false
			}

			return true
		})
	}

	return true, nil
}
