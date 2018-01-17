package maltmill

import (
	"io"

	"github.com/google/go-github/github"
)

type creator struct {
	writer    io.Writer
	overwrite bool
	outFile   string
	ghcli     *github.Client
}

func (cr *creator) run() error {
	return nil
}
