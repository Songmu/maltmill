class Kibelasync < Formula
  version '0.1.0'
  homepage 'https://github.com/Songmu/kibelasync'
  if OS.mac?
    url "https://github.com/Songmu/kibelasync/releases/download/v0.1.0/kibelasync_v0.1.0_darwin_amd64.zip"
    sha256 '758f07a1073c6924a4c09b167b413b915e623c342920092d655a7eb21cdd443b'
  end
  if OS.linux?
    url "https://github.com/Songmu/kibelasync/releases/download/v0.1.0/kibelasync_v0.1.0_linux_amd64.tar.gz"
    sha256 'bc92df3d0cb9aafd0a6449726ffdfd4dca348b1896d90a4ac4043561a59ec71d'
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
