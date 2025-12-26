package converter

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/vib795/epub2pdf/internal/epub"
)

// Options holds conversion options
type Options struct {
	PageSize    string  // A4, Letter, etc.
	Margin      float64 // Margin in inches
	Landscape   bool
	PrintBG     bool // Print background graphics
	Scale       float64
	Verbose     bool
}

// DefaultOptions returns sensible defaults
func DefaultOptions() Options {
	return Options{
		PageSize:  "A4",
		Margin:    0.5,
		Landscape: false,
		PrintBG:   true,
		Scale:     1.0,
		Verbose:   false,
	}
}

// Convert converts an EPUB book to PDF
func Convert(book *epub.Book, outputPath string, opts Options) error {
	html := book.ToHTML()

	// Create a temporary HTML file
	tmpFile, err := os.CreateTemp("", "epub2pdf-*.html")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(html); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write temp file: %w", err)
	}
	tmpFile.Close()

	// Create Chrome context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set timeout
	ctx, cancel = context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	// Navigate and print to PDF
	var pdfData []byte
	
	fileURL := "file://" + tmpFile.Name()
	
	if opts.Verbose {
		fmt.Printf("Converting HTML to PDF using headless Chrome...\n")
	}

	// Get page dimensions based on page size
	width, height := getPageDimensions(opts.PageSize)
	if opts.Landscape {
		width, height = height, width
	}

	err = chromedp.Run(ctx,
		chromedp.Navigate(fileURL),
		chromedp.WaitReady("body"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfData, _, err = page.PrintToPDF().
				WithPaperWidth(width).
				WithPaperHeight(height).
				WithMarginTop(opts.Margin).
				WithMarginBottom(opts.Margin).
				WithMarginLeft(opts.Margin).
				WithMarginRight(opts.Margin).
				WithPrintBackground(opts.PrintBG).
				WithScale(opts.Scale).
				Do(ctx)
			return err
		}),
	)

	if err != nil {
		return fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Write PDF to output file
	if err := os.WriteFile(outputPath, pdfData, 0644); err != nil {
		return fmt.Errorf("failed to write PDF: %w", err)
	}

	return nil
}

// getPageDimensions returns width and height in inches for common page sizes
func getPageDimensions(pageSize string) (float64, float64) {
	switch pageSize {
	case "Letter":
		return 8.5, 11
	case "Legal":
		return 8.5, 14
	case "Tabloid":
		return 11, 17
	case "A3":
		return 11.69, 16.54
	case "A4":
		return 8.27, 11.69
	case "A5":
		return 5.83, 8.27
	default:
		return 8.27, 11.69 // Default to A4
	}
}
