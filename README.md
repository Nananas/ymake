# Ymake

Ymake is a build tool analog to GNU make. Instead of using makefiles, ymake reads from a `Ymakefile`, written in [YAML](https://en.wikipedia.org/wiki/YAML).

## Command line 
Command line usage is as follows:

```
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
```

## Example Ymakefile
An example of a ymakefile is the one used to build ymake:

```
## YMAKEFILE for ymake itself

variables:
  gopath: echo $GOPATH
  version: echo 0.2

blocks:
  # Build
  default:
    pattern: "(*).go" # pattern not really needed though, will result in ymake.go only
    target: '{gopath}/bin/$1'
    deps:
      - "$1.go"
      - "ymakefile/*.go"
    cmd: go install -ldflags "-X main.VERSION={version}" github.com/nananas/$1 

```

### Variables
Variables can be specified using the `variables` key. Each variable has a name and a corresponding shell command. 
This command is executed and the result is replaced for every occurrence of a `{variable_name}` in any string in `cmd`, `pattern`, `deps` and `target`.

Default variables include:

- `WD`: the current working directory
- `HOME`: the home directory of the current user

### blocks
A block can contain the following:

- `cmd` <Either>: shell command to execute. (Can contain $)
- `pattern` <string>: simple pattern matching and capture. Wildcards `*?` can be used. Capturing is done using `(...)`.
- `deps` <Either>: dependency files. Needs `target`. (Can contain $)
- `target` <string>: file to check deps against. If it does not exist or a deps file is newer, the command will be executed. (Can contain $)
- `post` <Either>: any block names to run after this block.


Where <Either> is either a string or a list of strings:
```
# All viable:
deps: foobar
deps:
	- foo
	- bar
deps: [foo, bar]
```

Each `$1-$0` in `deps`, `target` and `cmd` will be replaced by the corresponding capture (`(...)`) from `pattern`.
A `$0` will be replaced by a list off all captures, joined by spaces.


# TODO

[o] error handling