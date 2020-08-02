package maltmill

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/sync/errgroup"
)

type cmdMaltmill struct {
	files     []string
	overwrite bool

	writer io.Writer

	ghcli *github.Client
}

var _ runner = (*cmdMaltmill)(nil)

func (mm *cmdMaltmill) run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, f := range mm.files {
		f := f
		eg.Go(func() error {
			return mm.processFile(ctx, f)
		})
	}
	return eg.Wait()
}

func (mm *cmdMaltmill) processFile(ctx context.Context, f string) error {
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
