package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := cmd[2]

	for key, e := range env {
		if e.NeedRemove {
			os.Unsetenv(key)
			continue
		}

		os.Setenv(key, e.Value)
	}

	c := exec.Command(command)
	c.Env = os.Environ()
	c.Args = cmd[2:]
	c.Stdout = os.Stdout

	if err := c.Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}

		return 1
	}

	return 0
}
