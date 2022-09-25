class Maltmill < Formula
  version '0.5.7'
  homepage 'https://github.com/Songmu/maltmill'

  on_macos do
    if Hardware::CPU.intel?
      url 'https://github.com/Songmu/maltmill/releases/download/v0.5.7/maltmill_v0.5.7_darwin_amd64.zip'
      sha256 '2af4eec3a80441e016514726efe630fac57ee30855b1f7c83f82c76e07f167e2'
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url 'https://github.com/Songmu/maltmill/releases/download/v0.5.7/maltmill_v0.5.7_linux_amd64.tar.gz'
      sha256 '2af4eec3a80441e016514726efe630fac57ee30855b1f7c83f82c76e07f167e4'
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
