package maltmill

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

type cmdNew struct {
	writer    io.Writer
	slug      string
	overwrite bool
	outFile   string
	ghcli     *github.Client
}

var _ runner = (*cmdNew)(nil)

var tmpl = `class {{.CapitalizedName}} < Formula
  version '{{.Version}}'
  homepage 'https://github.com/{{.Owner}}/{{.Repo}}'
  if OS.mac?
    url "{{.Downloads.darwin.URL}}"
    sha256 '{{.Downloads.darwin.SHA256}}'
  end
  if OS.linux?
    url "{{.Downloads.linux.URL}}"
    sha256 '{{.Downloads.linux.SHA256}}'
  end
  head 'https://github.com/{{.Owner}}/{{.Repo}}.git'

  head do
    depends_on 'go' => :build
  end

  def install
    if build.head?
      system 'make', 'build'
    end
    bin.install '{{.Name}}'
  end
end
`

type formulaData struct {
	Name, CapitalizedName string
	Version               string
	Owner, Repo           string
	Downloads             map[string]formulaDownload
}

type formulaDownload struct {
	URL    string
	SHA256 string
	OS     string
	Arch   string
}

var formulaTmpl = template.Must(template.New("formulaTmpl").Parse(tmpl))

var osNameRe = regexp.MustCompile("(darwin|linux)")

func getDownloads(assets []github.ReleaseAsset) ([]formulaDownload, error) {
	var downloads []formulaDownload
	for _, asset := range assets {
		u := asset.GetBrowserDownloadURL()
		fname := path.Base(u)
		if !strings.Contains(fname, "amd64") {
			continue
		}
		osName := osNameRe.FindString(fname)
		if osName == "" {
			continue
		}
		digest, err := getSHA256FromURL(u)
		if err != nil {
			return nil, err
		}
		downloads = append(downloads, formulaDownload{
			URL:    u,
			SHA256: digest,
			OS:     osName,
			Arch:   "amd64",
		})
	}
	if len(downloads) == 0 {
		return nil, errors.New("no assets found")
	}
	return downloads, nil
}

func (cr *cmdNew) run(ctx context.Context) (err error) {
	ownerAndRepo := strings.Split(cr.slug, "/")
	if len(ownerAndRepo) != 2 {
		return errors.Errorf("invalid slug: %s", cr.slug)
	}
	repoAndVer := strings.Split(ownerAndRepo[1], "@")
	var tag string
	if len(repoAndVer) > 1 {
		tag = repoAndVer[1]
	}
	nf := &formulaData{
		Owner:           ownerAndRepo[0],
		Repo:            repoAndVer[0],
		Name:            repoAndVer[0],
		CapitalizedName: strings.Replace(strings.Title(repoAndVer[0]), "-", "", -1),
	}
	rele, resp, err := func() (*github.RepositoryRelease, *github.Response, error) {
		if tag == "" {
			return cr.ghcli.Repositories.GetLatestRelease(context.Background(), nf.Owner, nf.Repo)
		}
		return cr.ghcli.Repositories.GetReleaseByTag(context.Background(), nf.Owner, nf.Repo, tag)
	}()
	if err != nil {
		return errors.Wrapf(err, "create new formula failed")
	}
	resp.Body.Close()

	ver, err := semver.NewVersion(rele.GetTagName())
	if err != nil {
		return errors.Wrapf(err, "invalid tag name: %s", rele.GetTagName())
	}
	nf.Version = fmt.Sprintf("%d.%d.%d", ver.Major(), ver.Minor(), ver.Patch())
	downloads, err := getDownloads(rele.Assets)
	if err != nil {
		return err
	}
	nf.Downloads = make(map[string]formulaDownload, len(downloads))
	for _, d := range downloads {
		nf.Downloads[d.OS] = d
	}
	var wtr = cr.writer
	if cr.overwrite || cr.outFile != "" {
		fname := cr.outFile
		if fname == "" {
			fname = nf.Name + ".rb"
		}
		defer func() {
			if err == nil {
				log.Printf("created %q\n", fname)
			}
		}()
		f, err := os.Create(fname)
		if err != nil {
			return err
		}
		defer f.Close()
		wtr = f
	}
	return formulaTmpl.Execute(wtr, nf)
}
