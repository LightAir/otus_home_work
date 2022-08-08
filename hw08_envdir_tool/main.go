package main

import (
	"fmt"
	"os"
)

func main() {
	environment, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Printf("read directory error: %v", err)
		os.Exit(1)
	}

	returnCode := RunCmd(os.Args, environment)

	os.Exit(returnCode)
}
