package main

import (
	"os"

	"github.com/Songmu/maltmill"
)

func main() {
	os.Exit(maltmill.Run(os.Args[1:]))
}
