package maltmill

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
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
	for _, f := range mm.files {
		mm.processFile(f)
	}
	return nil
}

var (
	nameReg = regexp.MustCompile(`(?m)^\s+name\s*=\s*['"](.*)["']`)
	verReg  = regexp.MustCompile(`(?m)^\s+version\s*['"](.*)["']`)
	homeReg = regexp.MustCompile(`(?m)^\s+homepage\s*['"](.*)["']`)
	urlReg  = regexp.MustCompile(`(?m)^\s+url\s*['"](.*)["']`)
	shaReg  = regexp.MustCompile(`(?m)\s+sha256\s*['"](.*)["']`)

	parseHomeReg = regexp.MustCompile(`^https://github.com/([^/]+)/([^/]+)`)
)

func (mm *maltmill) processFile(f string) error {
	_, err := getFormula(f)
	if err != nil {
		return err
	}
	return nil
}

func getFormula(f string) (*formula, error) {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	fo := &formula{fname: f}
	fo.content = string(b)

	if m := nameReg.FindStringSubmatch(fo.content); len(m) > 1 {
		fo.name = m[1]
	}
	if m := verReg.FindStringSubmatch(fo.content); len(m) < 2 {
		return nil, errors.New("no version detected")
	} else {
		fo.version = m[1]
	}
	if m := shaReg.FindStringSubmatch(fo.content); len(m) < 2 {
		return nil, errors.New("no sha256 detected")
	} else {
		fo.sha256 = m[1]
	}

	info := map[string]string{
		"name":    fo.name,
		"version": fo.version,
	}

	if m := homeReg.FindStringSubmatch(fo.content); len(m) < 2 {
		return nil, errors.New("no homepage detected")
	} else {
		h := m[1]
		fo.homepage, err = expandStr(h, info)
		if err != nil {
			return nil, err
		}
		if m := parseHomeReg.FindStringSubmatch(fo.homepage); len(m) < 3 {
			return nil, errors.Errorf("invalid homepage format: %s", fo.homepage)
		} else {
			fo.owner = m[1]
			fo.repo = m[2]
		}
	}

	if m := urlReg.FindStringSubmatch(fo.content); len(m) < 2 {
		return nil, errors.New("no url detected")
	} else {
		u := m[1]
		fo.url, err = expandStr(u, info)
		if err != nil {
			return nil, err
		}
	}

	return fo, nil
}

type formula struct {
	fname string

	content                              string
	name, version, homepage, url, sha256 string
	owner, repo                          string
}

func expandStr(str string, m map[string]string) (string, error) {
	for k, v := range m {
		reg, err := regexp.Compile(`#{` + k + `}`)
		if err != nil {
			return "", err
		}
		str = reg.ReplaceAllString(str, v)
	}
	return str, nil
}
