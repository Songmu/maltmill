package maltmill

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/v74/github"
	"github.com/pkg/errors"
)

type formula struct {
	fname string

	content       string
	name, version string
	owner, repo   string
	tagPrefix     string
}

var (
	nameReg = regexp.MustCompile(`(?m)^\s+name\s*=\s*['"](.*)["']`)
	verReg  = regexp.MustCompile(`(?m)(^\s+version\s*['"])(.*)(["'])`)
	urlReg  = regexp.MustCompile(`(?m)(^\s+url\s*['"])(.*)(["'])`)

	parseURLReg = regexp.MustCompile(`^https://[^/]*github.com/([^/]+)/([^/]+)`)
)

func newFormula(f string) (*formula, error) {
	b, err := os.ReadFile(f)
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

	info := map[string]string{
		"name":    fo.name,
		"version": fo.version,
	}

	m = urlReg.FindStringSubmatch(fo.content)
	if len(m) < 4 {
		return nil, errors.New("no url detected")
	}
	url := m[2]
	if strings.Contains(url, "#{version}") {
		url, err = expandStr(url, info)
		if err != nil {
			return nil, err
		}
	}

	m = parseURLReg.FindStringSubmatch(url)
	if len(m) < 3 {
		return nil, errors.Errorf("invalid url format: %s", url)
	}
	fo.owner = m[1]
	fo.repo = m[2]

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

func (fo *formula) update(ctx context.Context, ghcli *github.Client) (updated bool, err error) {
	origVer, err := semver.NewVersion(fo.version)
	if err != nil {
		return false, errors.Wrap(err, "invalid original version")
	}

	rele, newVer, err := fo.findLatestRelease(ctx, ghcli)
	if err != nil {
		return false, errors.Wrapf(err, "update formula failed: %s", fo.fname)
	}
	if !origVer.LessThan(newVer) {
		return false, nil
	}
	fromTag := fo.tagPrefix + fo.version
	fromRele, resp, err := ghcli.Repositories.GetReleaseByTag(ctx, fo.owner, fo.repo, fromTag)
	if err != nil {
		return false, errors.Wrapf(err, "update formula failed: %s", fo.fname)
	}
	resp.Body.Close()

	newVerStr := fmt.Sprintf("%d.%d.%d", newVer.Major(), newVer.Minor(), newVer.Patch())
	fromDownloads, err := getDownloads(fromRele.Assets)
	if err != nil {
		return false, errors.Wrapf(err, "update formula failed: %s", fo.fname)
	}
	downloads, err := getDownloads(rele.Assets)
	if err != nil {
		return false, errors.Wrapf(err, "update formula failed: %s", fo.fname)
	}

	fo.version = newVerStr
	fo.updateContent(fromDownloads, downloads)

	return true, nil
}

func parseTagVersionWithPrefix(tag, prefix string) (*semver.Version, error) {
	if !strings.HasPrefix(tag, prefix) {
		return nil, errors.Errorf("tag %q does not have prefix %q", tag, prefix)
	}
	version, err := semver.NewVersion(strings.TrimPrefix(tag, prefix))
	if err != nil {
		return nil, errors.Wrapf(err, "invalid tag name: %s", tag)
	}
	return version, nil
}

func selectLatestReleaseByPrefix(releases []*github.RepositoryRelease, prefix string) (*github.RepositoryRelease, *semver.Version) {
	var latest *github.RepositoryRelease
	var latestVer *semver.Version
	for _, rele := range releases {
		if rele.GetDraft() || rele.GetPrerelease() {
			continue
		}
		ver, err := parseTagVersionWithPrefix(rele.GetTagName(), prefix)
		if err != nil {
			continue
		}
		if latest == nil || latestVer.LessThan(ver) {
			latest = rele
			latestVer = ver
		}
	}
	return latest, latestVer
}

func findLatestReleaseByPrefix(ctx context.Context, ghcli *github.Client, owner, repo, prefix string) (*github.RepositoryRelease, *semver.Version, error) {
	var (
		latest    *github.RepositoryRelease
		latestVer *semver.Version
		opt       = &github.ListOptions{PerPage: 100}
	)

	for {
		releases, resp, err := ghcli.Repositories.ListReleases(ctx, owner, repo, opt)
		if resp != nil {
			resp.Body.Close()
		}
		if err != nil {
			return nil, nil, err
		}
		cand, candVer := selectLatestReleaseByPrefix(releases, prefix)
		if cand != nil && (latest == nil || latestVer.LessThan(candVer)) {
			latest = cand
			latestVer = candVer
		}
		if resp == nil || resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	if latest == nil {
		return nil, nil, errors.Errorf("no matching releases found for prefix %q", prefix)
	}
	return latest, latestVer, nil
}

func (fo *formula) findLatestRelease(ctx context.Context, ghcli *github.Client) (*github.RepositoryRelease, *semver.Version, error) {
	return findLatestReleaseByPrefix(ctx, ghcli, fo.owner, fo.repo, fo.tagPrefix)
}

func (fo *formula) updateContent(from, to []formulaDownload) {
	// Sort formulaDownloads by URL length in order to replace by longest match.
	sort.Slice(from, func(i, j int) bool {
		return len(from[i].URL) > len(from[j].URL)
	})
	sort.Slice(to, func(i, j int) bool {
		return len(to[i].URL) > len(to[j].URL)
	})
	var replacements []string
	for _, fromD := range from {
		for _, toD := range to {
			if fromD.Arch == toD.Arch && fromD.OS == toD.OS {
				replacements = append(replacements, fromD.URL, toD.URL, fromD.SHA256, toD.SHA256)
			}
		}
	}

	r := strings.NewReplacer(replacements...)
	fo.content = r.Replace(fo.content)
	fo.content = replaceOne(verReg, fo.content, fmt.Sprintf(`${1}%s${3}`, fo.version))
}

func replaceOne(reg *regexp.Regexp, str, replace string) string {
	replaced := false
	return reg.ReplaceAllStringFunc(str, func(match string) string {
		if replaced {
			return match
		}
		replaced = true
		return reg.ReplaceAllString(match, replace)
	})
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
