package maltmill

import (
	"flag"
	"log"
	"os"

	"github.com/google/go-github/github"
)

const (
	exitCodeOK = iota
	exitCodeErr
)

// Run the maltmill
func Run(args []string) int {
	err := (&cli{outStream: os.Stdout, errStream: os.Stderr}).run(args)
	if err != nil {
		if err == flag.ErrHelp {
			return exitCodeOK
		}
		log.Printf("[!!ERROR!!] %s\n", err)
		return exitCodeErr
	}
	return exitCodeOK
}

type maltmill struct {
	files []string

	ghcli *github.Client
}

func (mm *maltmill) run() error {
	return nil
}
