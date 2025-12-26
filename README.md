# epub2pdf

A fast, reliable command-line tool for converting EPUB files to PDF, written in Go.

## Features

- 📖 **Full EPUB Support** - Parses EPUB 2 and EPUB 3 formats
- 🎨 **Preserves Styling** - Maintains CSS styling and formatting
- 🖼️ **Image Embedding** - Embeds all images including covers as base64
- 📐 **Flexible Page Sizes** - A4, A5, A3, Letter, Legal, Tabloid
- 🔄 **Orientation Options** - Portrait or landscape mode
- ⚡ **Fast Conversion** - Uses headless Chrome for accurate rendering
- 🖥️ **Cross-Platform** - Works on Linux, macOS, and Windows

## Installation

### Homebrew (macOS/Linux)

```bash
# Add the tap
brew tap vib795/tap

# Install epub2pdf
brew install epub2pdf
```

### Go Install

```bash
go install github.com/vib795/epub2pdf@latest
```

### Download Binary

Download the latest release from the [Releases page](https://github.com/vib795/epub2pdf/releases).

### From Source

```bash
# Clone the repository
git clone https://github.com/vib795/epub2pdf.git
cd epub2pdf

# Build
make build

# Or install to GOPATH/bin
make install
```

## Prerequisites

- **Chrome/Chromium** installed for PDF rendering

If you installed via Homebrew on macOS and don't have Chrome:
```bash
brew install --cask chromium
```

## Usage

### Basic Conversion

```bash
# Convert EPUB to PDF (output: book.pdf)
epub2pdf book.epub

# Specify output path
epub2pdf book.epub -o output.pdf
```

### Options

```bash
epub2pdf [flags] <input.epub> [output.pdf]

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

# Verbose output
epub2pdf book.epub -v
```

### View EPUB Info

```bash
# Display EPUB metadata without converting
epub2pdf info book.epub
```

### Version

```bash
epub2pdf version
```

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

## How It Works

1. **Parse EPUB**: Opens the EPUB (ZIP archive), reads `container.xml` to find the OPF file
2. **Extract Content**: Parses the OPF manifest and spine to get chapters in reading order
3. **Build HTML**: Combines all chapters into a single styled HTML document
4. **Render PDF**: Uses headless Chrome (via chromedp) to render HTML to PDF

## Building Releases

```bash
# Build for all platforms
make release

# Output in releases/
# - epub2pdf-linux-amd64
# - epub2pdf-linux-arm64
# - epub2pdf-darwin-amd64
# - epub2pdf-darwin-arm64
# - epub2pdf-windows-amd64.exe
```

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
