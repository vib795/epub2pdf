#!/bin/bash
# scripts/update-formula.sh
# Updates the Homebrew formula with SHA256 checksums from the latest release

set -e

VERSION="${1:-}"
if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 1.0.0"
    exit 1
fi

REPO="vib795/epub2pdf"
FORMULA="Formula/epub2pdf.rb"

echo "📦 Updating formula for version $VERSION..."

# Download and calculate checksums
declare -A CHECKSUMS

for platform in darwin_amd64 darwin_arm64 linux_amd64 linux_arm64; do
    URL="https://github.com/$REPO/releases/download/v$VERSION/epub2pdf_${VERSION}_${platform}.tar.gz"
    echo "⬇️  Downloading $platform..."
    
    SHA=$(curl -sL "$URL" | shasum -a 256 | cut -d' ' -f1)
    CHECKSUMS[$platform]=$SHA
    echo "   SHA256: $SHA"
done

# Update formula
echo ""
echo "📝 Updating $FORMULA..."

# Create updated formula
cat > "$FORMULA" << EOF
class Epub2pdf < Formula
  desc "Fast CLI tool to convert EPUB files to PDF"
  homepage "https://github.com/$REPO"
  version "$VERSION"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/$REPO/releases/download/v$VERSION/epub2pdf_${VERSION}_darwin_amd64.tar.gz"
      sha256 "${CHECKSUMS[darwin_amd64]}"
    end
    on_arm do
      url "https://github.com/$REPO/releases/download/v$VERSION/epub2pdf_${VERSION}_darwin_arm64.tar.gz"
      sha256 "${CHECKSUMS[darwin_arm64]}"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/$REPO/releases/download/v$VERSION/epub2pdf_${VERSION}_linux_amd64.tar.gz"
      sha256 "${CHECKSUMS[linux_amd64]}"
    end
    on_arm do
      url "https://github.com/$REPO/releases/download/v$VERSION/epub2pdf_${VERSION}_linux_arm64.tar.gz"
      sha256 "${CHECKSUMS[linux_arm64]}"
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
EOF

echo "✅ Formula updated successfully!"
echo ""
echo "Next steps:"
echo "  1. Copy Formula/epub2pdf.rb to your homebrew-tap repo"
echo "  2. Commit and push the changes"
echo "  3. Users can then: brew upgrade epub2pdf"
