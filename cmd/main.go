package main

import (
	"fmt"

	"github.com/loureirovinicius/cleanup/cmd/cleaner"
)

func main() {
	err := cleaner.Run()
	if err != nil {
		panic(fmt.Errorf("there was an error running the cleaner: %v", err))
	}
}
