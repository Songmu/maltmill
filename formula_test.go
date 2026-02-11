package maltmill

import (
	"os"
	"reflect"
	"testing"

	"github.com/google/go-github/v74/github"
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

func TestParseTagVersionWithPrefix(t *testing.T) {
	testCases := []struct {
		name          string
		tag           string
		prefix        string
		expectVersion string
		expectErr     bool
	}{{
		name:          "v prefix",
		tag:           "v1.2.3",
		prefix:        "v",
		expectVersion: "1.2.3",
	}, {
		name:          "product prefix",
		tag:           "my-product-v0.8.1",
		prefix:        "my-product-v",
		expectVersion: "0.8.1",
	}, {
		name:      "too many version segments",
		tag:       "v1.2.3.4",
		prefix:    "v",
		expectErr: true,
	}, {
		name:      "prefix mismatch",
		tag:       "other-v1.2.3",
		prefix:    "my-product-v",
		expectErr: true,
	}, {
		name:      "invalid tag",
		tag:       "main",
		prefix:    "v",
		expectErr: true,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ver, err := parseTagVersionWithPrefix(tc.tag, tc.prefix)
			if tc.expectErr {
				if err == nil {
					t.Fatalf("error should not be nil for tag %q", tc.tag)
				}
				return
			}
			if err != nil {
				t.Fatalf("error should be nil but: %s", err)
			}
			if ver.String() != tc.expectVersion {
				t.Errorf("unexpected version. out=%s expect=%s", ver.String(), tc.expectVersion)
			}
		})
	}
}

func TestSelectLatestReleaseByPrefix(t *testing.T) {
	releases := []*github.RepositoryRelease{{
		TagName: github.String("my-product-v1.1.0"),
	}, {
		TagName: github.String("other-v9.9.9"),
	}, {
		TagName: github.String("my-product-v1.4.0"),
	}, {
		TagName:    github.String("my-product-v1.5.0"),
		Prerelease: github.Bool(true),
	}}

	rele, ver := selectLatestReleaseByPrefix(releases, "my-product-v")
	if rele == nil || ver == nil {
		t.Fatal("release and version should not be nil")
	}
	if rele.GetTagName() != "my-product-v1.4.0" {
		t.Errorf("unexpected release. out=%s expect=%s", rele.GetTagName(), "my-product-v1.4.0")
	}
	if ver.String() != "1.4.0" {
		t.Errorf("unexpected version. out=%s expect=%s", ver.String(), "1.4.0")
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
	}, {
		name:    "with arm64",
		fname:   "testdata/filt.rb",
		version: "0.8.2",
		fromDownloads: []formulaDownload{{
			URL:    "https://github.com/k1LoW/filt/releases/download/v0.8.0/filt_v0.8.0_darwin_arm64.zip",
			SHA256: "2923c5d872745728fc21ac0b5f8b892f987c4654f00591fa7555df6f0b52301f",
			OS:     "darwin",
			Arch:   "arm64",
		}, {
			URL:    "https://github.com/k1LoW/filt/releases/download/v0.8.0/filt_v0.8.0_darwin_amd64.zip",
			SHA256: "137b15437a15569141b12c02d05317e1c0b9db567ea7ef4d88e11cce0db0d1fd",
			OS:     "darwin",
			Arch:   "amd64",
		}, {
			URL:    "https://github.com/k1LoW/filt/releases/download/v0.8.0/filt_v0.8.0_linux_arm64.tar.gz",
			SHA256: "5424589e650c66de2c238473ed53bf1c7ab99f9dfb91eeffc5437bd9d397de89",
			OS:     "linux",
			Arch:   "arm64",
		}, {
			URL:    "https://github.com/k1LoW/filt/releases/download/v0.8.0/filt_v0.8.0_linux_amd64.tar.gz",
			SHA256: "af1d99a7e0507b256eb1e88f36a5b507cca35952c13b0e00afa90a40d95418cb",
			OS:     "linux",
			Arch:   "amd64",
		}},
		downloads: []formulaDownload{{
			URL:    "https://github.com/k1LoW/filt/releases/download/v0.8.2/filt_v0.8.2_darwin_arm64.zip",
			SHA256: "bf72725b231cc01a72b16c4ef4bada598189529f70ef4b7c0a5fbccf18c0489e",
			OS:     "darwin",
			Arch:   "arm64",
		}, {
			URL:    "https://github.com/k1LoW/filt/releases/download/v0.8.2/filt_v0.8.2_darwin_amd64.zip",
			SHA256: "2579fcabf8ca89ed278500abe8489cf0d082b985a5b2e9065cfd45b3bd8ae324",
			OS:     "darwin",
			Arch:   "amd64",
		}, {
			URL:    "https://github.com/k1LoW/filt/releases/download/v0.8.2/filt_v0.8.2_linux_arm64.tar.gz",
			SHA256: "f7e3e2951420f181a71e8d641087cd0eef8588eb9e017f4f4f8040067ed8624e",
			OS:     "linux",
			Arch:   "arm64",
		}, {
			URL:    "https://github.com/k1LoW/filt/releases/download/v0.8.2/filt_v0.8.2_linux_amd64.tar.gz",
			SHA256: "9c052ef5749960e5a848713fb143de75958fbbe95a277f778a741c3063e0d182",
			OS:     "linux",
			Arch:   "amd64",
		}},
		expectFile: "testdata/filt_update.rb",
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
