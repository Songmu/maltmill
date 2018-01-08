class Goxz < Formula
  name = 'goxz'
  version '0.1.0'
  homepage "https://github.com/Songmu/#{name}"
  url "https://github.com/Songmu/#{name}/releases/download/v#{version}/#{name}_v#{version}_darwin_amd64.zip"
  sha256 '1449899f3e49615b4cbb17493a2f63b88a7489bb4ffb0b0b7a9992e6508cab38'

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
