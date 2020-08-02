package main

import (
	"context"
	"os"

	"github.com/Songmu/maltmill"
)

func main() {
	os.Exit(maltmill.Run(context.Background(), os.Args[1:], os.Stdout, os.Stderr))
}
