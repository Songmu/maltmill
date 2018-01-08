package maltmill

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	"github.com/pkg/errors"
	gitconfig "github.com/tcnksm/go-gitconfig"
)

const (
	envGitHubToken = "GITHUB_TOKEN"
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
	var token string
	fs.StringVar(&token, "token", os.Getenv(envGitHubToken), "")

	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}
	mm.files = fs.Args()
	if len(mm.files) < 1 {
		return nil, errors.New("no formula files are specified")
	}

	if token == "" {
		token, _ = gitconfig.GithubToken()
	}
	mm.ghcli = newGithubClient(token)

	return mm, nil
}


