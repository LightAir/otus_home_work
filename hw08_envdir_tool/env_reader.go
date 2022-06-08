package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	environment := make(Environment)

	for _, f := range files {
		file, err := os.Open(dir + "/" + f.Name())
		if err != nil {
			return nil, err
		}
		defer file.Close()

		reader := bufio.NewReader(file)
		line, _ := reader.ReadString('\n')
		line = strings.TrimRight(line, "\t\n\r ")
		line = strings.ReplaceAll(line, string(rune(0x00)), "\n")

		env := EnvValue{
			Value:      line,
			NeedRemove: len(line) == 0,
		}

		name := strings.ReplaceAll(f.Name(), "=", "")

		environment[name] = env
	}

	return environment, nil
}
