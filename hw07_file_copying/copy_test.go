package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const tempFileName = "tmp"

func TestCopy(t *testing.T) {
	t.Run("unsupported file", func(t *testing.T) {
		err := Copy("/dev/urandom", tempFileName, 0, 0)
		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "Error should be: %v, got: %v", ErrUnsupportedFile, err)
	})

	t.Run("offset exceeds file size", func(t *testing.T) {
		err := Copy("testdata/input.txt", tempFileName, 10000, 0)

		isOffsetError := errors.Is(err, ErrOffsetExceedsFileSize)
		require.Truef(t, isOffsetError, "Should be: %v, got: %v", ErrOffsetExceedsFileSize, err)
	})

	t.Run("copy jpeg file", func(t *testing.T) {
		err := Copy("testdata/pic.jpeg", tempFileName, 0, 0)
		require.Nilf(t, err, "Expected: nil, got %v", err)
	})

	t.Run("Source file not found", func(t *testing.T) {
		err := Copy("testdata/not_found.jpeg", "", 0, 0)
		errorMsg := "copy from testdata/not_found.jpeg: open testdata/not_found.jpeg: no such file or directory"

		require.Equalf(t, err.Error(), errorMsg, "Error should be: %v, got: %v", errorMsg, err)
	})

	t.Run("file not permission on write", func(t *testing.T) {
		err := Copy("testdata/input.txt", "/not_writable", 0, 0)
		errorMsg := "file creation error /not_writable: open /not_writable: permission denied"

		require.NotNilf(t, err, "Expected: %v, got nil", errorMsg)
		require.Equalf(t, err.Error(), errorMsg, "Error should be: %v, got: %v", errorMsg, err)
	})

	t.Run("limit less file size", func(t *testing.T) {
		err := Copy("testdata/pic.jpeg", tempFileName, 0, 10)
		require.Nilf(t, err, "Expected: nil, got %v", err)
	})

	t.Run("copy error", func(t *testing.T) {
		err := Copy("/", tempFileName, 0, 10)
		errorMsg := "bytes copying error from / to tmp: read /: is a directory"

		require.NotNilf(t, err, "Expected: %v, got nil", errorMsg)
		require.Equalf(t, err.Error(), errorMsg, "Error should be: %v, got: %v", errorMsg, err)
	})

	t.Run("eof", func(t *testing.T) {
		err := Copy("testdata/input.txt", tempFileName, 0, 0)
		require.Nilf(t, err, "Expected: nil, got %v", err)
	})

	t.Run("offset greater zero", func(t *testing.T) {
		err := Copy("testdata/input.txt", tempFileName, 50, 0)
		require.Nilf(t, err, "Expected: nil, got %v", err)
	})

	t.Run("big offset", func(t *testing.T) {
		err := Copy("/", tempFileName, 100000, 0)
		errorMsg := "offset exceeds file size"

		require.NotNilf(t, err, "Expected: %v, got nil", errorMsg)
		require.Equalf(t, err.Error(), errorMsg, "Error should be: %v, got: %v", errorMsg, err)
	})

	tearDown()
}

func tearDown() {
	_ = os.Remove(tempFileName)
}
