package maltmill

import (
	"os"
	"reflect"
	"testing"
)

func TestNewFormula(t *testing.T) {
	fname := "testdata/goxz.rb"
	fo, err := newFormula(fname)
	if err != nil {
		t.Errorf("err should be nil but: %s", err)
	}
	fo.content = ""

	expect := formula{
		fname:   fname,
		name:    "goxz",
		version: "0.1.0",
		owner:   "Songmu",
		repo:    "goxz",
	}

	if !reflect.DeepEqual(*fo, expect) {
		t.Errorf("failed to getFormula.\n   out: %#v\nexpect: %#v", *fo, expect)
	}
}

func TestUpdateContent(t *testing.T) {
	testCases := []struct {
		name          string
		fname         string
		version       string
		fromDownloads []formulaDownload
		downloads     []formulaDownload
		expectFile    string
	}{{
		name:    "template url",
		fname:   "testdata/goxz.rb",
		version: "0.2.1",
		fromDownloads: []formulaDownload{{
			SHA256: "1449899f3e49615b4cbb17493a2f63b88a7489bb4ffb0b0b7a9992e6508cab38",
			URL:    "https://github.com/Songmu/goxz/releases/download/v0.1.0/goxz_v0.1.0_darwin_amd64.zip",
			OS:     "darwin",
			Arch:   "amd64",
		}},
		downloads: []formulaDownload{{
			SHA256: "11112222",
			URL:    "https://github.com/Songmu/goxz/releases/download/v0.2.1/goxz_v0.2.1_darwin_amd64.zip",
			OS:     "darwin",
			Arch:   "amd64",
		}},
		expectFile: "testdata/goxz_update.rb",
	}, {
		name:    "raw url",
		fname:   "testdata/goxz2.rb",
		version: "0.3.3",
		fromDownloads: []formulaDownload{{
			SHA256: "1449899f3e49615b4cbb17493a2f63b88a7489bb4ffb0b0b7a9992e6508cab38",
			URL:    "https://github.com/Songmu/goxz/releases/download/v0.1.0/goxz_v0.1.0_darwin_amd64.zip",
			OS:     "darwin",
			Arch:   "amd64",
		}},
		downloads: []formulaDownload{{
			SHA256: "11113333",
			URL:    "https://github.com/Songmu/goxz/releases/download/v0.3.3/goxz_v0.3.3_darwin_amd64.zip",
			OS:     "darwin",
			Arch:   "amd64",
		}},
		expectFile: "testdata/goxz2_update.rb",
	}, {
		name:    "with linux",
		fname:   "testdata/kibelasync.rb",
		version: "0.1.1",
		fromDownloads: []formulaDownload{{
			SHA256: "758f07a1073c6924a4c09b167b413b915e623c342920092d655a7eb21cdd443b",
			URL:    "https://github.com/Songmu/kibelasync/releases/download/v0.1.0/kibelasync_v0.1.0_darwin_amd64.zip",
			OS:     "darwin",
			Arch:   "amd64",
		}, {
			SHA256: "bc92df3d0cb9aafd0a6449726ffdfd4dca348b1896d90a4ac4043561a59ec71d",
			URL:    "https://github.com/Songmu/kibelasync/releases/download/v0.1.0/kibelasync_v0.1.0_linux_amd64.tar.gz",
			OS:     "linux",
			Arch:   "amd64",
		}},
		downloads: []formulaDownload{{
			SHA256: "11118888",
			URL:    "https://github.com/Songmu/kibelasync/releases/download/v0.1.1/kibelasync_v0.1.1_darwin_amd64.zip",
			OS:     "darwin",
			Arch:   "amd64",
		}, {
			SHA256: "11119999",
			URL:    "https://github.com/Songmu/kibelasync/releases/download/v0.1.1/kibelasync_v0.1.1_linux_amd64.tar.gz",
			OS:     "linux",
			Arch:   "amd64",
		}},
		expectFile: "testdata/kibelasync_update.rb",
	}, {
		name:    "partial matching url",
		fname:   "testdata/ecspresso.rb",
		version: "0.18.0",
		fromDownloads: []formulaDownload{{
			URL:    "https://github.com/kayac/ecspresso/releases/download/v0.17.3/ecspresso-v0.17.3-darwin-amd64",
			SHA256: "1ac91503dcf2e7883b9df0d2a5b54c7c9a49e2d7b78f73286cfc19bb6ad44778",
			OS:     "darwin",
			Arch:   "amd64",
		}, {
			URL:    "https://github.com/kayac/ecspresso/releases/download/v0.17.3/ecspresso-v0.17.3-darwin-amd64.zip",
			SHA256: "34684ce9b841eec0d30c809081bbf8269c1b9456301282fb54cf907d6687743d",
			OS:     "darwin",
			Arch:   "amd64",
		}},
		downloads: []formulaDownload{{
			URL:    "https://github.com/kayac/ecspresso/releases/download/v0.18.0/ecspresso-v0.18.0-darwin-amd64",
			SHA256: "44f7f90acf75ee38a18b50f5ff90de5b6d5ef8d3b639cf0942244a8a699e6aef",
			OS:     "darwin",
			Arch:   "amd64",
		}, {
			URL:    "https://github.com/kayac/ecspresso/releases/download/v0.18.0/ecspresso-v0.18.0-darwin-amd64.zip",
			SHA256: "7fee5a7c401afd84e9b099aba37bea41572fd29731ddc3e4afe4bea4b3470a36",
			OS:     "darwin",
			Arch:   "amd64",
		}},
		expectFile: "testdata/ecspresso_update.rb",
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fo, err := newFormula(tc.fname)
			if err != nil {
				t.Errorf("err should be nil but: %s", err)
			}
			fo.version = tc.version
			fo.updateContent(tc.fromDownloads, tc.downloads)

			b, _ := os.ReadFile(tc.expectFile)
			expect := string(b)

			if fo.content != expect {
				t.Errorf("something went wrong.\n  out=%s\nexpect=%s", fo.content, expect)
			}
		})
	}
}

func TestGetSHA256FromURL(t *testing.T) {
	out, err := getSHA256FromURL("https://github.com/Songmu/go-sandbox/releases/download/v0.1.0/go-sandbox_v0.1.0_darwin_amd64.zip")

	if err != nil {
		t.Errorf("error should be nil but: %s", err)
	}

	expect := "4b4a1e064c2c3534edadca4f532c712367fd0f22148ae8f994850a0407635c0a"
	if out != expect {
		t.Errorf("unexpected sha256.\n   out: %s\nexpect: %s", out, expect)
	}
}
