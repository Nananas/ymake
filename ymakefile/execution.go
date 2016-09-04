package ymakefile

import (
	"os"
	"os/exec"
	"strings"
)

// Execution to stdout and stderr
//
func ExecuteStd(command string, stdin string) error {
	cmd := exec.Command("/bin/sh", "-c", command)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin + "\n")
	}

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// Execution to string
//
func ExecuteCapture(command string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Stdin = os.Stdin

	b, err := cmd.CombinedOutput()
	if err != nil {
		return string(b), err
	}

	r := string(b)

	// trim whitespace
	r = strings.Trim(r, "\n")
	if err != nil {
		return r, err
	}

	return r, nil
}
