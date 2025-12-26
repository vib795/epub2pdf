.PHONY: build install clean test lint run

# Build variables
BINARY_NAME=epub2pdf
VERSION=1.0.0
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X github.com/vib795/epub2pdf/cmd.Version=$(VERSION) \
                  -X github.com/vib795/epub2pdf/cmd.BuildDate=$(BUILD_DATE) \
                  -X github.com/vib795/epub2pdf/cmd.GitCommit=$(GIT_COMMIT)"

# Default target
all: build

# Build the binary
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) .
	@echo "✅ Build complete: bin/$(BINARY_NAME)"

# Install to GOPATH/bin
install:
	@echo "📦 Installing $(BINARY_NAME)..."
	go install $(LDFLAGS) .
	@echo "✅ Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -rf bin/
	go clean
	@echo "✅ Clean complete"

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Run linter
lint:
	@echo "🔍 Running linter..."
	golangci-lint run

# Download dependencies
deps:
	@echo "📥 Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies updated"

# Build for multiple platforms
release:
	@echo "🚀 Building releases..."
	@mkdir -p releases
	
	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o releases/$(BINARY_NAME)-linux-amd64 .
	
	# Linux ARM64
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o releases/$(BINARY_NAME)-linux-arm64 .
	
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o releases/$(BINARY_NAME)-darwin-amd64 .
	
	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o releases/$(BINARY_NAME)-darwin-arm64 .
	
	# Windows AMD64
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o releases/$(BINARY_NAME)-windows-amd64.exe .
	
	@echo "✅ Releases built in releases/"

# Run example
run:
	@go run . --help

# Show help
help:
	@echo "epub2pdf Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build    - Build the binary"
	@echo "  make install  - Install to GOPATH/bin"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make test     - Run tests"
	@echo "  make lint     - Run linter"
	@echo "  make deps     - Download dependencies"
	@echo "  make release  - Build for all platforms"
	@echo "  make run      - Run with --help"
	@echo "  make help     - Show this help"

# Update Homebrew formula (after creating a release)
formula:
	@echo "📦 Updating Homebrew formula..."
	@./scripts/update-formula.sh $(VERSION)

# Create a new release tag
tag:
	@if [ -z "$(VERSION)" ]; then echo "Usage: make tag VERSION=1.0.0"; exit 1; fi
	@echo "🏷️  Creating tag v$(VERSION)..."
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	git push origin v$(VERSION)
	@echo "✅ Tag v$(VERSION) pushed. GitHub Actions will create the release."
