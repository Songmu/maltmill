package maltmill

import (
	"io/ioutil"
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
	}{
		{
			name:    "template url",
			fname:   "testdata/goxz.rb",
			version: "0.2.1",
			fromDownloads: []formulaDownload{{
				SHA256: "1449899f3e49615b4cbb17493a2f63b88a7489bb4ffb0b0b7a9992e6508cab38",
				URL:    "https://github.com/Songmu/goxz/releases/download/v0.1.0/goxz_v0.1.0_darwin_amd64.zip",
			}},
			downloads: []formulaDownload{{
				SHA256: "11112222",
				URL:    "https://github.com/Songmu/goxz/releases/download/v0.2.1/goxz_v0.2.1_darwin_amd64.zip",
			}},
			expectFile: "testdata/goxz_update.rb",
		},
		{
			name:    "raw url",
			fname:   "testdata/goxz2.rb",
			version: "0.3.3",
			fromDownloads: []formulaDownload{{
				SHA256: "1449899f3e49615b4cbb17493a2f63b88a7489bb4ffb0b0b7a9992e6508cab38",
				URL:    "https://github.com/Songmu/goxz/releases/download/v0.1.0/goxz_v0.1.0_darwin_amd64.zip",
			}},
			downloads: []formulaDownload{{
				SHA256: "11113333",
				URL:    "https://github.com/Songmu/goxz/releases/download/v0.3.3/goxz_v0.3.3_darwin_amd64.zip",
			}},
			expectFile: "testdata/goxz2_update.rb",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fo, err := newFormula(tc.fname)
			if err != nil {
				t.Errorf("err should be nil but: %s", err)
			}
			fo.version = tc.version
			fo.updateContent(tc.fromDownloads, tc.downloads)

			b, _ := ioutil.ReadFile(tc.expectFile)
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
