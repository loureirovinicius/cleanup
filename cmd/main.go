package main

import (
	"fmt"
	"os"

	"github.com/loureirovinicius/cleanup/cmd/cleaner"
)

func main() {
	err := cleaner.Run()
	if err != nil {
		fmt.Printf("there was an error running the cleaner: %v\n", err)
		os.Exit(1)
	}
}
