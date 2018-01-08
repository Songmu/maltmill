package maltmill

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/sync/errgroup"
)

const (
	exitCodeOK = iota
	exitCodeErr
)

// Run the maltmill
func Run(args []string) int {
	err := (&cli{outStream: os.Stdout, errStream: os.Stderr}).run(args)
	if err != nil {
		if err == flag.ErrHelp || err == errOpt {
			return exitCodeOK
		}
		log.Printf("[!!ERROR!!] %s\n", err)
		return exitCodeErr
	}
	return exitCodeOK
}

type maltmill struct {
	files     []string
	overwrite bool

	writer io.Writer

	ghcli *github.Client
}

func (mm *maltmill) run() error {
	eg := errgroup.Group{}
	for _, f := range mm.files {
		f := f
		eg.Go(func() error {
			return mm.processFile(f)
		})
	}
	return eg.Wait()
}

func (mm *maltmill) processFile(f string) error {
	fo, err := newFormula(f)
	if err != nil {
		return err
	}
	updated, err := fo.update(mm.ghcli)
	if err != nil {
		return err
	}
	if mm.overwrite && !updated {
		return nil
	}

	var w io.Writer = mm.writer
	if mm.overwrite {
		f, err := os.Create(f)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}
	_, err = fmt.Fprint(w, fo.content)
	return err
}
