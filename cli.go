package maltmill

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	"github.com/Songmu/ghselfupdate"
	"github.com/pkg/errors"
)

const envGitHubToken = "GITHUB_TOKEN"

var errOpt = errors.New("special option requested")

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

func (cl *cli) parseArgs(args []string) (runner, error) {
	mm := &maltmill{writer: cl.outStream}
	fs := flag.NewFlagSet("maltmill", flag.ContinueOnError)
	fs.SetOutput(cl.errStream)
	fs.Usage = func() {
		fs.SetOutput(cl.outStream)
		defer fs.SetOutput(cl.errStream)
		fmt.Fprintf(cl.outStream, `maltmill - Update homebrew third party formula

Version: %s (rev: %s/%s)

Synopsis:
    %% maltmill -w [formula-files.rb]

Options:
`, version, revision, runtime.Version())
		fs.PrintDefaults()
		fmt.Fprintf(cl.outStream, `
Commands:
    new            create new formula
	self-update    binary self update
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
		return cl.parseCreatorArgs(append(newArgs, restArgs[1:]...))
	case "self-update":
		return &updator{}, nil
	default:
		mm.files = restArgs
		mm.ghcli = newGithubClient(token)
		return mm, nil
	}
}

type updator struct {
}

func (upd *updator) run() error {
	return ghselfupdate.Do(version)
}

func (cl *cli) parseCreatorArgs(args []string) (runner, error) {
	cr := &creator{writer: cl.outStream}
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
	cr.ghcli = newGithubClient(token)
	return cr, nil
}
