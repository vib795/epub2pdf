class Epub2pdf < Formula
  desc "Fast CLI tool to convert EPUB files to PDF"
  homepage "https://github.com/vib795/epub2pdf"
  version "1.0.0"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/vib795/epub2pdf/releases/download/v1.0.0/epub2pdf_1.0.0_darwin_amd64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256_DARWIN_AMD64"
    end
    on_arm do
      url "https://github.com/vib795/epub2pdf/releases/download/v1.0.0/epub2pdf_1.0.0_darwin_arm64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256_DARWIN_ARM64"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/vib795/epub2pdf/releases/download/v1.0.0/epub2pdf_1.0.0_linux_amd64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256_LINUX_AMD64"
    end
    on_arm do
      url "https://github.com/vib795/epub2pdf/releases/download/v1.0.0/epub2pdf_1.0.0_linux_arm64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256_LINUX_ARM64"
    end
  end

  depends_on "chromium" => :optional

  def install
    bin.install "epub2pdf"
  end

  def caveats
    <<~EOS
      epub2pdf requires Chrome or Chromium for PDF rendering.
      
      If you don't have Chrome installed, you can install Chromium:
        brew install --cask chromium
      
      Or download Chrome from: https://www.google.com/chrome/
    EOS
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/epub2pdf version")
  end
end
