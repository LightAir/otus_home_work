package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		env, err := ReadDir("testdata/env_custom")
		require.NotNil(t, env, "env can't be nil")
		require.Nilf(t, err, "err must be nil")

		expectedEnvs := Environment{
			"BAR":          EnvValue{Value: "bar", NeedRemove: false},
			"EMPTY":        EnvValue{Value: "", NeedRemove: true},
			"FOO":          EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
			"REMOVEEQUALS": EnvValue{Value: "remove equals from name", NeedRemove: false},
			"TAB_CHECK":    EnvValue{Value: "\"next tabs\"", NeedRemove: false},
			"UNSET":        EnvValue{Value: "", NeedRemove: true},
		}

		require.EqualValuesf(t, expectedEnvs, env, "Expected: %v, got %v", "asd", env)
	})

	t.Run("bad dir", func(t *testing.T) {
		env, err := ReadDir("folder_not_found")
		errMsg := "open folder_not_found: no such file or directory"

		require.Nilf(t, env, "Expected: nil, got %v", env)
		require.NotNilf(t, err, "Expected: %v, got nil", errMsg)
		require.EqualErrorf(t, err, errMsg, "Error should be: %v, got: %v", errMsg, err)
	})

	t.Run("folder contains bad files", func(t *testing.T) {
		env, err := ReadDir("/dev/block")
		errMsg := "permission denied"

		require.Nilf(t, env, "Expected: nil, got %v", env)
		require.NotNilf(t, err, "Expected: %v, got nil", errMsg)
		require.ErrorContains(t, err, errMsg, "Error should contains: %v, got: %v", errMsg, err)
	})
}
