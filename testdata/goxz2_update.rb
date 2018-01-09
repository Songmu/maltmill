class Goxz < Formula
  name = 'goxz'
  version '0.3.3'
  homepage "https://github.com/Songmu/#{name}"
  url "https://github.com/Songmu/goxz/releases/download/v0.3.3/goxz_v0.3.3_darwin_amd64.zip"
  sha256 '11113333'

  head do
    url "https://github.com/Songmu/#{name}.git"
    depends_on 'go' => :build
  end

  def install
    if build.head?
      gobin = buildpath/"bin"
      ENV.update({
        'GOPATH' => buildpath,
        'PATH'   => "#{gobin}:#{ENV['PATH']}",
      })
      mkdir_p buildpath/'src/github.com/Songmu'
      ln_s buildpath, buildpath/"src/github.com/Songmu/#{name}"
      system 'go', 'get', '-d', '-v', './...'
      system 'go', 'build', "./cmd/#{name}"
    end

    bin.install name
  end
end
