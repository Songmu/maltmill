package maltmill

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

type creator struct {
	writer    io.Writer
	slug      string
	overwrite bool
	outFile   string
	ghcli     *github.Client
}

var tmpl = `class {{.CapitalizedName}} < Formula
  version '{{.Version}}'
  homepage 'https://github.com/{{.Owner}}/{{.Repo}}'
  url "{{.URL}}"
  sha256 '{{.SHA256}}'
  head 'https://github.com/{{.Owner}}/{{.Repo}}.git'

  def install
    bin.install '{{.Name}}'
  end
end
`

type NewFormula struct {
	Name, CapitalizedName string
	Version               string
	Owner, Repo           string
	SHA256, URL           string
}

var formulaTmpl = template.Must(template.New("formulaTmpl").Parse(tmpl))

func (cr *creator) run() error {
	ownerAndRepo := strings.Split(cr.slug, "/")
	if len(ownerAndRepo) != 2 {
		return errors.Errorf("invalid slug: %s", cr.slug)
	}
	nf := &NewFormula{
		Owner:           ownerAndRepo[0],
		Repo:            ownerAndRepo[1],
		Name:            ownerAndRepo[1],
		CapitalizedName: strings.Title(ownerAndRepo[1]),
	}
	rele, resp, err := cr.ghcli.Repositories.GetLatestRelease(context.Background(), nf.Owner, nf.Repo)
	if err != nil {
		return errors.Wrapf(err, "create new formula failed")
	}
	resp.Body.Close()

	ver, err := semver.NewVersion(rele.GetTagName())
	if err != nil {
		return errors.Wrapf(err, "invalid tag name: %s", rele.GetTagName())
	}
	nf.Version = fmt.Sprintf("%d.%d.%d", ver.Major(), ver.Minor(), ver.Patch())
	nf.URL, err = func() (string, error) {
		for _, asset := range rele.Assets {
			u := asset.GetBrowserDownloadURL()
			fname := path.Base(u)
			if strings.Contains(fname, "amd64") &&
				strings.Contains(fname, "darwin") {
				return u, nil
			}
		}
		return "", errors.New("no assets found from latest release")
	}()
	if err != nil {
		return err
	}
	nf.SHA256, err = getSHA256FromURL(nf.URL)
	if err != nil {
		return errors.Wrapf(err, "faild to create new formula")
	}

	var wtr io.Writer = cr.writer
	if cr.overwrite {
		fname := cr.outFile
		if fname == "" {
			fname = nf.Name + ".rb"
		}
		f, err := os.Create(fname)
		if err != nil {
			return err
		}
		defer f.Close()
		wtr = f
	}
	return formulaTmpl.Execute(wtr, nf)
}
