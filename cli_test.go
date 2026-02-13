package maltmill

import (
	"bytes"
	"context"
	"io"
	"testing"
)

func TestParseArgsTagPrefix(t *testing.T) {
	cl := &cli{
		outStream: io.Discard,
		errStream: io.Discard,
	}
	ctx := context.Background()

	rnr, err := cl.parseArgs(ctx, []string{"-tag-prefix", "my-product-v", "testdata/goxz.rb"})
	if err != nil {
		t.Fatalf("parseArgs should not fail but: %s", err)
	}
	mm, ok := rnr.(*cmdMaltmill)
	if !ok {
		t.Fatalf("runner should be *cmdMaltmill but: %T", rnr)
	}
	if mm.tagPrefix != "my-product-v" {
		t.Errorf("unexpected tagPrefix. out=%s expect=%s", mm.tagPrefix, "my-product-v")
	}
}

func TestParseArgsTagPrefixForNew(t *testing.T) {
	cl := &cli{
		outStream: io.Discard,
		errStream: io.Discard,
	}
	ctx := context.Background()

	rnr, err := cl.parseArgs(ctx, []string{"-tag-prefix", "tool-v", "new", "Songmu/maltmill"})
	if err != nil {
		t.Fatalf("parseArgs should not fail but: %s", err)
	}
	cr, ok := rnr.(*cmdNew)
	if !ok {
		t.Fatalf("runner should be *cmdNew but: %T", rnr)
	}
	if cr.tagPrefix != "tool-v" {
		t.Errorf("unexpected tagPrefix. out=%s expect=%s", cr.tagPrefix, "tool-v")
	}
}

func TestNew(t *testing.T) {
	out := &bytes.Buffer{}
	cl := &cli{
		outStream: out,
		errStream: io.Discard,
	}
	ctx := context.Background()
	rnr, err := cl.parseArgs(ctx, []string{"new", "Songmu/maltmill@v0.0.1"})
	if err != nil {
		t.Errorf("error should be nil on parsing but: %s", err)
	}
	err = rnr.run(ctx)
	if err != nil {
		t.Errorf("error should be nil but: %s", err)
	}
	expect := `class Maltmill < Formula
  desc 'create and update Homebrew thrid party Formulae'
  version '0.0.1'
  homepage 'https://github.com/Songmu/maltmill'

  on_macos do
    if Hardware::CPU.intel?
      url 'https://github.com/Songmu/maltmill/releases/download/v0.0.1/maltmill_v0.0.1_darwin_amd64.zip'
      sha256 '7433c1c1e48eb05601bbf91b0ffb76f5f298773a2b87a584088a9a7562562969'
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url 'https://github.com/Songmu/maltmill/releases/download/v0.0.1/maltmill_v0.0.1_linux_amd64.tar.gz'
      sha256 'c77dbf0053ca718172b506886f8cb55deab859d2a50598aa014959bae758e7b4'
    end
  end

  head do
    url 'https://github.com/Songmu/maltmill.git'
    depends_on 'go' => :build
  end

  def install
    if build.head?
      system 'make', 'build'
    end
    bin.install 'maltmill'
  end
end
`

	if out.String() != expect {
		t.Errorf("result not expected.\n  out: %s\nexpect: %s", out.String(), expect)
	}
}
