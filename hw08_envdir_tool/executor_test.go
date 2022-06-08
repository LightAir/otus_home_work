package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	envs := Environment{
		"BAR":   EnvValue{Value: "bar", NeedRemove: false},
		"EMPTY": EnvValue{Value: "", NeedRemove: true},
		"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
		"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
		"UNSET": EnvValue{Value: "", NeedRemove: true},
	}

	t.Run("bad command", func(t *testing.T) {
		cmd := []string{"", "", "/bin/"}

		returnCode := RunCmd(cmd, envs)

		require.Equalf(t, 1, returnCode, "Expected: 1, got %v", returnCode)
	})

	t.Run("good command", func(t *testing.T) {
		cmd := []string{"", "", "/bin/ls"}

		returnCode := RunCmd(cmd, envs)

		require.Equalf(t, 0, returnCode, "Expected: 0, got %v", returnCode)
	})

	t.Run("bad command with error return code", func(t *testing.T) {
		cmd := []string{"", "", "/bin/ls", "/not_found_dir/"}

		returnCode := RunCmd(cmd, envs)

		require.Equalf(t, 2, returnCode, "Expected: 2, got %v", returnCode)
	})
}
