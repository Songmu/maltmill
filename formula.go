package maltmill

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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
	verReg  = regexp.MustCompile(`(?m)(^\s+version\s*['"])(.*)(["'])`)
	homeReg = regexp.MustCompile(`(?m)^\s+homepage\s*['"](.*)["']`)
	urlReg  = regexp.MustCompile(`(?m)^\s+url\s*['"](.*)["']`)
	shaReg  = regexp.MustCompile(`(?m)(\s+sha256\s*['"])(.*)(["'])`)

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
	m := verReg.FindStringSubmatch(fo.content)
	if len(m) < 4 {
		return nil, errors.New("no version detected")
	}
	fo.version = m[2]

	m = shaReg.FindStringSubmatch(fo.content)
	if len(m) < 4 {
		return nil, errors.New("no sha256 detected")
	}
	fo.sha256 = m[2]

	info := map[string]string{
		"name":    fo.name,
		"version": fo.version,
	}

	m = homeReg.FindStringSubmatch(fo.content)
	if len(m) < 2 {
		return nil, errors.New("no homepage detected")
	}
	h := m[1]
	fo.homepage, err = expandStr(h, info)
	if err != nil {
		return nil, err
	}
	m = parseHomeReg.FindStringSubmatch(fo.homepage)
	if len(m) < 3 {
		return nil, errors.Errorf("invalid homepage format: %s", fo.homepage)
	}
	fo.owner = m[1]
	fo.repo = m[2]

	m = urlReg.FindStringSubmatch(fo.content)
	if len(m) < 2 {
		return nil, errors.New("no url detected")
	}
	fo.urlTmpl = m[1]
	fo.url, err = expandStr(fo.urlTmpl, info)
	if err != nil {
		return nil, err
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
		return false, errors.Wrapf(err, "update formula failed: %s", fo.fname)
	}
	resp.Body.Close()

	newVer, err := semver.NewVersion(rele.GetTagName())
	if err != nil {
		return false, errors.Wrapf(err, "invalid original version. formula: %s", fo.fname)
	}
	if !origVer.LessThan(newVer) {
		return false, nil
	}

	newVerStr := fmt.Sprintf("%d.%d.%d", newVer.Major(), newVer.Minor(), newVer.Patch())
	newURL, err := expandStr(fo.urlTmpl, map[string]string{
		"name":    fo.name,
		"version": newVerStr,
	})
	if err != nil {
		return false, errors.Wrapf(err, "faild to upload formula: %s", fo.fname)
	}
	newSHA256, err := getSHA256FromURL(newURL)
	if err != nil {
		return false, errors.Wrapf(err, "faild to upload formula: %s", fo.fname)
	}
	fo.version = newVerStr
	fo.url = newURL
	fo.sha256 = newSHA256
	fo.updateContent()

	return true, nil
}

// update version and sha256
func (fo *formula) updateContent() {
	fo.content = verReg.ReplaceAllString(fo.content, fmt.Sprintf(`${1}%s${3}`, fo.version))
	fo.content = shaReg.ReplaceAllString(fo.content, fmt.Sprintf(`${1}%s${3}`, fo.sha256))
}

func getSHA256FromURL(u string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", fmt.Sprintf("maltmill/%s", version))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrapf(err, "getSHA256 failed while request to url: %s", u)
	}
	defer resp.Body.Close()

	h := sha256.New()
	if _, err := io.Copy(h, resp.Body); err != nil {
		return "", errors.Wrapf(err, "getSHA256 failed while reading response. url: %s", u)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
