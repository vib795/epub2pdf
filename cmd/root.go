package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vib795/epub2pdf/internal/converter"
	"github.com/vib795/epub2pdf/internal/epub"
)

var (
	// Flags
	outputPath string
	pageSize   string
	margin     float64
	landscape  bool
	noBG       bool
	scale      float64
	verbose    bool
)

var rootCmd = &cobra.Command{
	Use:   "epub2pdf <input.epub> [output.pdf]",
	Short: "Convert EPUB files to PDF",
	Long: `epub2pdf is a command-line tool for converting EPUB ebooks to PDF format.

It parses the EPUB structure, extracts chapters in reading order,
and renders them to a beautifully formatted PDF document.

Examples:
  epub2pdf book.epub                    # Output: book.pdf
  epub2pdf book.epub -o output.pdf      # Specify output path
  epub2pdf book.epub --page-size Letter # Use US Letter size
  epub2pdf book.epub --landscape        # Landscape orientation
  epub2pdf book.epub -v                 # Verbose output`,
	Args: cobra.MinimumNArgs(1),
	RunE: runConvert,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output PDF path (default: input name with .pdf extension)")
	rootCmd.Flags().StringVarP(&pageSize, "page-size", "p", "A4", "Page size: A4, A5, Letter, Legal, Tabloid")
	rootCmd.Flags().Float64VarP(&margin, "margin", "m", 0.5, "Page margin in inches")
	rootCmd.Flags().BoolVarP(&landscape, "landscape", "l", false, "Use landscape orientation")
	rootCmd.Flags().BoolVar(&noBG, "no-background", false, "Don't print background graphics")
	rootCmd.Flags().Float64VarP(&scale, "scale", "s", 1.0, "Scale factor (0.1 - 2.0)")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}

func runConvert(cmd *cobra.Command, args []string) error {
	inputPath := args[0]

	// Validate input file
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return fmt.Errorf("input file not found: %s", inputPath)
	}

	if !strings.HasSuffix(strings.ToLower(inputPath), ".epub") {
		return fmt.Errorf("input file must be an EPUB file")
	}

	// Determine output path
	output := outputPath
	if output == "" {
		base := strings.TrimSuffix(inputPath, filepath.Ext(inputPath))
		output = base + ".pdf"
	}

	// Validate scale
	if scale < 0.1 || scale > 2.0 {
		return fmt.Errorf("scale must be between 0.1 and 2.0")
	}

	// Validate page size
	validSizes := map[string]bool{
		"A4": true, "A5": true, "A3": true,
		"Letter": true, "Legal": true, "Tabloid": true,
	}
	if !validSizes[pageSize] {
		return fmt.Errorf("invalid page size: %s (valid: A4, A5, A3, Letter, Legal, Tabloid)", pageSize)
	}

	if verbose {
		fmt.Printf("📖 Input:  %s\n", inputPath)
		fmt.Printf("📄 Output: %s\n", output)
		fmt.Printf("📐 Page:   %s", pageSize)
		if landscape {
			fmt.Print(" (landscape)")
		}
		fmt.Println()
	}

	// Parse EPUB
	if verbose {
		fmt.Println("🔍 Parsing EPUB...")
	}

	book, err := epub.Parse(inputPath)
	if err != nil {
		return fmt.Errorf("failed to parse EPUB: %w", err)
	}

	if verbose {
		fmt.Printf("📚 Title:    %s\n", book.Title)
		fmt.Printf("✍️  Author:   %s\n", book.Author)
		fmt.Printf("📑 Chapters: %d\n", len(book.Chapters))
	}

	// Convert to PDF
	if verbose {
		fmt.Println("🔄 Converting to PDF...")
	}

	opts := converter.Options{
		PageSize:  pageSize,
		Margin:    margin,
		Landscape: landscape,
		PrintBG:   !noBG,
		Scale:     scale,
		Verbose:   verbose,
	}

	if err := converter.Convert(book, output, opts); err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}

	// Get file size
	info, err := os.Stat(output)
	if err == nil {
		size := formatFileSize(info.Size())
		fmt.Printf("✅ Successfully created %s (%s)\n", output, size)
	} else {
		fmt.Printf("✅ Successfully created %s\n", output)
	}

	return nil
}

func formatFileSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	default:
		return fmt.Sprintf("%d bytes", size)
	}
}
