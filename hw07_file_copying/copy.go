package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const progressBarLength = 35

type counter struct {
	progress int64
	total    int64
}

func (c *counter) Write(p []byte) (n int, err error) {
	c.progress += int64(len(p))
	c.PrintProgressBar()

	return len(p), nil
}

func (c *counter) PrintProgressBar() {
	percent := int((c.progress * 100) / c.total)

	pbProgress := (percent * progressBarLength) / 100
	pbProgressStr := strings.Repeat("▒", pbProgress)
	bpLeftString := strings.Repeat("-", progressBarLength-pbProgress)

	fmt.Printf("\r%d / %d [%s%s] %d%%", c.progress, c.total, pbProgressStr, bpLeftString, percent)
}

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	from, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("copy from %s: %w", fromPath, err)
	}
	defer from.Close()

	fileInfo, err := from.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file information from %s: %w", fromPath, err)
	}

	fileSize := fileInfo.Size()
	if limit == 0 {
		limit = fileSize
	}

	if fileSize == 0 {
		return ErrUnsupportedFile
	}

	if offset >= fileSize {
		return ErrOffsetExceedsFileSize
	}

	_, err = from.Seek(offset, 0)
	if err != nil {
		return fmt.Errorf("set offset error for %s: %w", fromPath, err)
	}

	to, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("file creation error %s: %w", toPath, err)
	}
	defer to.Close()

	total := fileSize - offset
	if limit > 0 && limit+offset < fileSize {
		total = limit
	}

	reader := bufio.NewReader(from)
	writer := bufio.NewWriter(to)

	teeReader := io.TeeReader(reader, &counter{total: total})

	_, err = io.CopyN(writer, teeReader, limit)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return fmt.Errorf("bytes copying error from %s to %s: %w", fromPath, toPath, err)
		}
		return nil
	}

	return nil
}
