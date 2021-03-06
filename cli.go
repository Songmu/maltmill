package maltmill

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	"github.com/pkg/errors"
)

const envGitHubToken = "GITHUB_TOKEN"

var errOpt = errors.New("special option requested")

type cli struct {
	outStream, errStream io.Writer
}

func (cl *cli) run(ctx context.Context, args []string) error {
	log.SetOutput(cl.errStream)
	log.SetPrefix("[maltmill] ")
	log.SetFlags(0)

	mm, err := cl.parseArgs(ctx, args)
	if err != nil {
		return err
	}
	return mm.run(ctx)
}

func (cl *cli) parseArgs(ctx context.Context, args []string) (runner, error) {
	mm := &cmdMaltmill{writer: cl.outStream}
	fs := flag.NewFlagSet("maltmill", flag.ContinueOnError)
	fs.SetOutput(cl.errStream)
	fs.Usage = func() {
		fs.SetOutput(cl.outStream)
		defer fs.SetOutput(cl.errStream)
		fmt.Fprintf(cl.outStream, `maltmill - Update or create homebrew third party formulae

Version: %s (rev: %s/%s)

Synopsis:
    %% maltmill -w [formula-files.rb]

Options:
`, version, revision, runtime.Version())
		fs.PrintDefaults()
		fmt.Fprintf(cl.outStream, `
Commands:
    new            create new formula
`)
	}
	var token string
	fs.StringVar(&token, "token", os.Getenv(envGitHubToken), "github `token")
	fs.BoolVar(&mm.overwrite, "w", false, "write result to (source) file instead of stdout")

	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}

	restArgs := fs.Args()
	if len(restArgs) < 1 {
		return nil, errors.New("no formula files or sub command are specified")
	}
	switch restArgs[0] {
	case "new":
		newArgs := []string{}
		if token != "" {
			newArgs = append(newArgs, "-token", token)
		}
		if mm.overwrite {
			newArgs = append(newArgs, "-w")
		}
		return cl.parseCmdNewArgs(ctx, append(newArgs, restArgs[1:]...))
	default:
		mm.files = restArgs
		mm.ghcli = newGithubClient(ctx, token)
		return mm, nil
	}
}

func (cl *cli) parseCmdNewArgs(ctx context.Context, args []string) (runner, error) {
	cr := &cmdNew{writer: cl.outStream}
	fs := flag.NewFlagSet("maltmill", flag.ContinueOnError)
	fs.SetOutput(cl.errStream)
	fs.Usage = func() {
		fs.SetOutput(cl.outStream)
		defer fs.SetOutput(cl.errStream)
		fmt.Fprintf(cl.outStream, `maltmill new - create new formula

Synopsis:
    %% maltmill new -w Songmu/ghg

Options:
`)
		fs.PrintDefaults()
	}
	var token string
	fs.StringVar(&token, "token", os.Getenv(envGitHubToken), "github `token`")
	fs.BoolVar(&cr.overwrite, "w", false, "write result to (source) file instead of stdout")
	fs.StringVar(&cr.outFile, "o", "", "`file` to output")

	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}
	if len(fs.Args()) < 1 {
		return nil, errors.New("githut repository isn't specified")
	}
	cr.slug = fs.Arg(0)
	cr.ghcli = newGithubClient(ctx, token)
	return cr, nil
}
