package maltmill

import (
	"context"
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

type formula struct {
	fname string

	content                              string
	urlTmpl                              string
	name, version, homepage, url, sha256 string
	owner, repo                          string
}

var (
	nameReg = regexp.MustCompile(`(?m)^\s+name\s*=\s*['"](.*)["']`)
	verReg  = regexp.MustCompile(`(?m)^\s+version\s*['"](.*)["']`)
	homeReg = regexp.MustCompile(`(?m)^\s+homepage\s*['"](.*)["']`)
	urlReg  = regexp.MustCompile(`(?m)^\s+url\s*['"](.*)["']`)
	shaReg  = regexp.MustCompile(`(?m)\s+sha256\s*['"](.*)["']`)

	parseHomeReg = regexp.MustCompile(`^https://github.com/([^/]+)/([^/]+)`)
)

func newFormula(f string) (*formula, error) {
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
		fo.urlTmpl = m[1]
		fo.url, err = expandStr(fo.urlTmpl, info)
		if err != nil {
			return nil, err
		}
	}

	return fo, nil
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

func (fo *formula) update(ghcli *github.Client) (updated bool, err error) {
	origVer, err := semver.NewVersion(fo.version)
	if err != nil {
		return false, errors.Wrap(err, "invalid original version")
	}

	rele, resp, err := ghcli.Repositories.GetLatestRelease(context.Background(), fo.owner, fo.repo)
	if err != nil {
		return false, err
	}
	resp.Body.Close()

	newVer, err := semver.NewVersion(rele.GetTagName())
	if err != nil {
		return false, errors.Wrap(err, "invalid original version")
	}
	if !origVer.LessThan(newVer) {
		return false, nil
	}

	newVerStr := fmt.Sprintf("%d.%d.%d", newVer.Major(), newVer.Minor(), newVer.Patch())
	fo.version = newVerStr

	return false, nil
}
