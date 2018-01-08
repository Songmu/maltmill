package maltmill

import (
	"reflect"
	"testing"
)

func TestGetFormula(t *testing.T) {
	fname := "testdata/goxz.rb"
	fo, err := getFormula(fname)
	if err != nil {
		t.Errorf("err should be nil but: %s", err)
	}
	fo.content = ""

	expect := formula{
		fname:    fname,
		name:     "goxz",
		version:  "0.1.0",
		homepage: "https://github.com/Songmu/goxz",
		urlTmpl:  "https://github.com/Songmu/#{name}/releases/download/v#{version}/#{name}_v#{version}_darwin_amd64.zip",
		url:      "https://github.com/Songmu/goxz/releases/download/v0.1.0/goxz_v0.1.0_darwin_amd64.zip",
		sha256:   "1449899f3e49615b4cbb17493a2f63b88a7489bb4ffb0b0b7a9992e6508cab38",
		owner:    "Songmu",
		repo:     "goxz",
	}

	if !reflect.DeepEqual(*fo, expect) {
		t.Errorf("failed to getFormula.\n   out: %#v\nexpect: %#v", *fo, expect)
	}
}
