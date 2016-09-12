package ymakefile

import (
	"os"
	"os/exec"
	"strings"
)

var DEFAULT_SHELL = "/bin/sh"

// Executes to stdout and stderr
// Shell is a string, and has to have the -c option to specify the command
// Default is DEFAULT_SHELL
//
// the errorchan channel will receive the error if one occured, and resend it make it bubble to other processes
// the donechan channel is the communication with the parent process, sending a message if the command was completed
//
func ExecuteStdParallel(command, stdin, shell string, errorchan chan error, donechan chan bool) {

	if PRINTONLY {
		// return without executing
		donechan <- true
	}

	if shell == "" {
		shell = DEFAULT_SHELL
	}

	cmd := exec.Command(shell, "-c", command)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin + "\n")
	}

	// communication between the execute and state processes below
	//
	done := make(chan bool)

	// command execute process
	//
	go func() {
		err := cmd.Run()
		if err != nil {
			// signal other cmd state processes that an error occured
			//
			errorchan <- err
		}

		// signal parent
		//
		done <- true

	}()

	// command state process
	//
	go func() {
		select {
		case err := <-errorchan:
			cmd.Process.Kill()
			errorchan <- err
		case <-done:
			//
		}

		// exit normally
		donechan <- true
	}()

}

// Execution to stdout and stderr
// Shell is a string, and has to have the -c option to specify the command
// Default is DEFAULT_SHELL
//
func ExecuteStd(command, stdin, shell string) error {

	if PRINTONLY {
		return nil
	}

	if shell == "" {
		shell = DEFAULT_SHELL
	}

	cmd := exec.Command(shell, "-c", command)

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
func ExecuteCapture(command, shell string) (string, error) {

	if shell == "" {
		shell = DEFAULT_SHELL
	}

	cmd := exec.Command(shell, "-c", command)
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
