# epub2pdf

A fast, reliable command-line tool for converting EPUB files to PDF, written in Go.

<p align="center">
  <img src="demo.gif" alt="epub2pdf demo" width="600">
</p>

## Features

- 📖 **Full EPUB Support** - Parses EPUB 2 and EPUB 3 formats
- 🎨 **Preserves Styling** - Maintains CSS styling and formatting
- 🖼️ **Image Embedding** - Embeds all images including covers as base64
- 📐 **Flexible Page Sizes** - A4, A5, A3, Letter, Legal, Tabloid
- 🔄 **Orientation Options** - Portrait or landscape mode
- ⚡ **Fast Conversion** - Uses headless Chrome for accurate rendering
- 🖥️ **Cross-Platform** - Works on Linux, macOS, and Windows

## Installation

### Homebrew (macOS/Linux) — Recommended

```bash
brew tap vib795/tap
brew install epub2pdf
```

To update to the latest version:
```bash
brew update
brew upgrade epub2pdf
```

### Go Install

```bash
go install github.com/vib795/epub2pdf@latest
```

### Download Binary

Download the latest binary for your platform from the [Releases page](https://github.com/vib795/epub2pdf/releases):

| Platform | Download |
|----------|----------|
| macOS (Apple Silicon) | `epub2pdf_x.x.x_darwin_arm64.tar.gz` |
| macOS (Intel) | `epub2pdf_x.x.x_darwin_amd64.tar.gz` |
| Linux (x64) | `epub2pdf_x.x.x_linux_amd64.tar.gz` |
| Linux (ARM64) | `epub2pdf_x.x.x_linux_arm64.tar.gz` |
| Windows (x64) | `epub2pdf_x.x.x_windows_amd64.zip` |

Extract and move to your PATH:
```bash
tar -xzf epub2pdf_*.tar.gz
sudo mv epub2pdf /usr/local/bin/
```

### From Source

```bash
git clone https://github.com/vib795/epub2pdf.git
cd epub2pdf
make build
sudo mv bin/epub2pdf /usr/local/bin/
```

## Prerequisites

epub2pdf requires **Chrome or Chromium** for PDF rendering.

- **macOS**: Chrome is usually pre-installed, or run `brew install --cask chromium`
- **Linux**: `sudo apt install chromium-browser` or `sudo dnf install chromium`
- **Windows**: Download from [google.com/chrome](https://www.google.com/chrome/)

## Usage

### Basic Conversion

```bash
# Convert EPUB to PDF (outputs book.pdf)
epub2pdf book.epub

# Specify output path
epub2pdf book.epub -o output.pdf
```

### Options

```
epub2pdf [flags] <input.epub>

Flags:
  -o, --output string      Output PDF path (default: input name with .pdf)
  -p, --page-size string   Page size: A4, A5, Letter, Legal, Tabloid (default "A4")
  -m, --margin float       Page margin in inches (default 0.5)
  -l, --landscape          Use landscape orientation
      --no-background      Don't print background graphics
  -s, --scale float        Scale factor 0.1-2.0 (default 1.0)
  -v, --verbose            Verbose output
  -h, --help               Help for epub2pdf
```

### Examples

```bash
# Convert with US Letter size
epub2pdf book.epub --page-size Letter

# Landscape orientation with custom margins
epub2pdf book.epub -l -m 0.75

# Scale down for smaller file size
epub2pdf book.epub -s 0.8

# Verbose output to see progress
epub2pdf book.epub -v
```

### View EPUB Info

```bash
# Display EPUB metadata without converting
epub2pdf info book.epub
```

### Check Version

```bash
epub2pdf version
```

## How It Works

1. **Parse EPUB**: Opens the EPUB (ZIP archive), reads `container.xml` to find the OPF file
2. **Extract Content**: Parses the OPF manifest and spine to get chapters in reading order
3. **Embed Images**: Converts all images to base64 data URIs for self-contained HTML
4. **Build HTML**: Combines all chapters into a single styled HTML document
5. **Render PDF**: Uses headless Chrome (via chromedp) to render HTML to PDF

## Project Structure

```
epub2pdf/
├── main.go                     # Entry point
├── cmd/
│   ├── root.go                 # Main convert command
│   ├── info.go                 # Info subcommand
│   └── version.go              # Version subcommand
├── internal/
│   ├── epub/
│   │   └── parser.go           # EPUB parsing logic
│   └── converter/
│       └── converter.go        # HTML to PDF conversion
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Troubleshooting

### "Chrome not found" error
Make sure Chrome or Chromium is installed and accessible in your PATH.

### Images not appearing in PDF
Ensure your EPUB file contains valid image references. Run with `-v` for verbose output.

### PDF is too large
Use the scale option to reduce size: `epub2pdf book.epub -s 0.8`

## License

MIT License - see [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request