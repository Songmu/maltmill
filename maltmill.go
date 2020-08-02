package maltmill

import (
	"context"
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
func Run(ctx context.Context, args []string, outStream, errStream io.Writer) int {
	err := (&cli{outStream: outStream, errStream: errStream}).run(ctx, args)
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

func (mm *maltmill) run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, f := range mm.files {
		f := f
		eg.Go(func() error {
			return mm.processFile(ctx, f)
		})
	}
	return eg.Wait()
}

func (mm *maltmill) processFile(ctx context.Context, f string) error {
	fo, err := newFormula(f)
	if err != nil {
		return err
	}
	updated, err := fo.update(ctx, mm.ghcli)
	if err != nil {
		return err
	}
	if mm.overwrite && !updated {
		return nil
	}

	w := mm.writer
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
