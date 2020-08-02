class Kibelasync < Formula
  version '0.1.1'
  homepage 'https://github.com/Songmu/kibelasync'
  if OS.mac?
    url "https://github.com/Songmu/kibelasync/releases/download/v0.1.1/kibelasync_v0.1.1_darwin_amd64.zip"
    sha256 '11118888'
  end
  if OS.linux?
    url "https://github.com/Songmu/kibelasync/releases/download/v0.1.1/kibelasync_v0.1.1_linux_amd64.tar.gz"
    sha256 '11119999'
  end
  head 'https://github.com/Songmu/kibelasync.git'

  head do
    depends_on 'go' => :build
  end

  def install
    if build.head?
      system 'make', 'build'
    end
    bin.install 'kibelasync'
  end
end
