package maltmill

import (
	"flag"
	"log"
	"os"
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
}

func (mm *maltmill) run() error {
	return nil
}
