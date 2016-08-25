package ymakefile

import (
	"os"
	"os/exec"
	"strings"
)

// Execution to stdout and stderr
//
func ExecuteStd(command string) error {
	cmd := exec.Command("/bin/sh", "-c", command)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

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

	b, err := cmd.CombinedOutput()
	r := string(b)

	// trim whitespace
	r = strings.Trim(r, "\n")
	if err != nil {
		return r, err
	}

	return r, nil
}
