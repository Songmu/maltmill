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

	"github.com/google/go-github/v74/github"
	"github.com/pkg/errors"
)

type cmdNew struct {
	writer    io.Writer
	slug      string
	overwrite bool
	outFile   string
	tagPrefix string
	ghcli     *github.Client
}

var _ runner = (*cmdNew)(nil)

var tmpl = `class {{.CapitalizedName}} < Formula
{{- if .Desc }}
  desc '{{.Desc | escapeSingleQuotes}}'
{{- end }}
  version '{{.Version}}'
  homepage 'https://github.com/{{.Owner}}/{{.Repo}}'
{{ if or (ne .Downloads.DarwinAmd64 nil) (ne .Downloads.DarwinArm64 nil) }}
  on_macos do
{{- if .Downloads.DarwinArm64 }}
    if Hardware::CPU.arm?
      url '{{.Downloads.DarwinArm64.URL}}'
      sha256 '{{.Downloads.DarwinArm64.SHA256}}'
    end
{{- end }}
{{- if .Downloads.DarwinAmd64 }}
    if Hardware::CPU.intel?
      url '{{.Downloads.DarwinAmd64.URL}}'
      sha256 '{{.Downloads.DarwinAmd64.SHA256}}'
    end
{{- end }}
  end
{{ end -}}
{{ if or (ne .Downloads.LinuxAmd64 nil) (ne .Downloads.LinuxArm64 nil) }}
  on_linux do
{{- if .Downloads.LinuxArm64 }}
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url '{{.Downloads.LinuxArm64.URL}}'
      sha256 '{{.Downloads.LinuxArm64.SHA256}}'
    end
{{- end }}
{{- if .Downloads.LinuxAmd64 }}
    if Hardware::CPU.intel?
      url '{{.Downloads.LinuxAmd64.URL}}'
      sha256 '{{.Downloads.LinuxAmd64.SHA256}}'
    end
{{- end }}
  end
{{ end }}
  head do
    url 'https://github.com/{{.Owner}}/{{.Repo}}.git'
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
	Desc                  string
	Downloads             formulaDataDownloads
}

type formulaDataDownloads struct {
	DarwinAmd64 *formulaDownload
	DarwinArm64 *formulaDownload
	LinuxAmd64  *formulaDownload
	LinuxArm64  *formulaDownload
}

type formulaDownload struct {
	URL    string
	SHA256 string
	OS     string
	Arch   string
}

func escapeSingleQuotes(in string) string {
	if in == "" {
		return ""
	}
	// Escape backslashes first to avoid double escaping.
	in = strings.ReplaceAll(in, "\\", "\\\\")
	in = strings.ReplaceAll(in, "'", "\\'")
	return in
}

var formulaTmpl = template.Must(
	template.New("formulaTmpl").Funcs(template.FuncMap{
		"escapeSingleQuotes": escapeSingleQuotes,
	}).Parse(tmpl),
)

var osNameRe = regexp.MustCompile("(darwin|linux)")

func getDownloads(assets []*github.ReleaseAsset) ([]formulaDownload, error) {
	var downloads []formulaDownload
	for _, asset := range assets {
		u := asset.GetBrowserDownloadURL()
		fname := path.Base(u)
		arch, ok := detectArch(fname)
		if !ok {
			continue
		}
		osName := osNameRe.FindString(fname)
		if osName == "" {
			continue
		}
		var digest string
		if asset.Digest != nil && strings.HasPrefix(*asset.Digest, "sha256:") {
			// if it is not a sha256 digest, we need to calculate from URL
			digest = strings.TrimPrefix(*asset.Digest, "sha256:")
		} else {
			var err error
			digest, err = getSHA256FromURL(u)
			if err != nil {
				return nil, err
			}
		}
		downloads = append(downloads, formulaDownload{
			URL:    u,
			SHA256: digest,
			OS:     osName,
			Arch:   arch,
		})
	}
	if len(downloads) == 0 {
		return nil, errors.New("no assets found")
	}
	return downloads, nil
}

func detectArch(in string) (string, bool) {
	archs := []string{"amd64", "arm64"}
	for _, a := range archs {
		if strings.Contains(in, a) {
			return a, true
		}
	}
	return "", false
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
		Downloads:       formulaDataDownloads{},
	}
	repo, resp, err := cr.ghcli.Repositories.Get(ctx, nf.Owner, nf.Repo)
	if err != nil {
		return errors.Wrapf(err, "create new formula failed")
	}
	nf.Desc = repo.GetDescription()
	resp.Body.Close()
	var rele *github.RepositoryRelease
	if tag != "" {
		rele, resp, err = cr.ghcli.Repositories.GetReleaseByTag(ctx, nf.Owner, nf.Repo, tag)
		if err != nil {
			return errors.Wrapf(err, "create new formula failed")
		}
		resp.Body.Close()
	} else {
		rele, _, err = findLatestReleaseByPrefix(ctx, cr.ghcli, nf.Owner, nf.Repo, cr.tagPrefix)
		if err != nil {
			return errors.Wrapf(err, "create new formula failed")
		}
	}

	ver, err := parseTagVersionWithPrefix(rele.GetTagName(), cr.tagPrefix)
	if err != nil {
		return errors.Wrapf(err, "invalid tag name: %s", rele.GetTagName())
	}
	nf.Version = fmt.Sprintf("%d.%d.%d", ver.Major(), ver.Minor(), ver.Patch())
	downloads, err := getDownloads(rele.Assets)
	if err != nil {
		return err
	}
	for _, d := range downloads {
		dd := d
		switch {
		case d.OS == "darwin" && d.Arch == "amd64":
			nf.Downloads.DarwinAmd64 = &dd
		case d.OS == "darwin" && d.Arch == "arm64":
			nf.Downloads.DarwinArm64 = &dd
		case d.OS == "linux" && d.Arch == "amd64":
			nf.Downloads.LinuxAmd64 = &dd
		case d.OS == "linux" && d.Arch == "arm64":
			nf.Downloads.LinuxArm64 = &dd
		}
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
