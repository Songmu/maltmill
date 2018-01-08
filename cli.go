package maltmill

import (
	"flag"
	"fmt"
	"io"
	"log"
	"runtime"
)

type cli struct {
	outStream, errStream io.Writer
}

func (cl *cli) run(args []string) error {
	log.SetOutput(cl.errStream)
	log.SetPrefix("[maltmill] ")
	log.SetFlags(0)

	mm, err := cl.parseArgs(args)
	if err != nil {
		return err
	}
	return mm.run()
}

func (cl *cli) parseArgs(args []string) (*maltmill, error) {
	mm := &maltmill{}
	fs := flag.NewFlagSet("maltmill", flag.ContinueOnError)
	fs.SetOutput(cl.errStream)
	fs.Usage = func() {
		fs.SetOutput(cl.outStream)
		defer fs.SetOutput(cl.errStream)
		fmt.Fprintf(cl.outStream, `maltmill - Update homebrew third party formula

Version: %s (rev: %s/%s)

Synopsis:
    %% maltmill [formula-files.rb]

Options:
`, version, revision, runtime.Version())
		fs.PrintDefaults()
	}

	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}
	return mm, nil
}
