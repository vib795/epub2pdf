package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vib795/epub2pdf/internal/epub"
)

var infoCmd = &cobra.Command{
	Use:   "info <input.epub>",
	Short: "Display EPUB metadata and structure",
	Long: `Display information about an EPUB file without converting it.

Shows the book title, author, number of chapters, and chapter list.

Example:
  epub2pdf info book.epub`,
	Args: cobra.ExactArgs(1),
	RunE: runInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func runInfo(cmd *cobra.Command, args []string) error {
	inputPath := args[0]

	// Validate input file
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return fmt.Errorf("input file not found: %s", inputPath)
	}

	if !strings.HasSuffix(strings.ToLower(inputPath), ".epub") {
		return fmt.Errorf("input file must be an EPUB file")
	}

	// Parse EPUB
	book, err := epub.Parse(inputPath)
	if err != nil {
		return fmt.Errorf("failed to parse EPUB: %w", err)
	}

	// Display info
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║                      EPUB Information                       ║")
	fmt.Println("╠════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ File:     %-49s ║\n", truncate(inputPath, 49))
	fmt.Printf("║ Title:    %-49s ║\n", truncate(book.Title, 49))
	fmt.Printf("║ Author:   %-49s ║\n", truncate(book.Author, 49))
	fmt.Printf("║ Chapters: %-49d ║\n", len(book.Chapters))
	fmt.Printf("║ CSS:      %-49d ║\n", len(book.CSS))
	fmt.Println("╠════════════════════════════════════════════════════════════╣")
	fmt.Println("║                      Chapter List                           ║")
	fmt.Println("╠════════════════════════════════════════════════════════════╣")

	for i, chapter := range book.Chapters {
		if i >= 20 {
			remaining := len(book.Chapters) - 20
			fmt.Printf("║   ... and %d more chapters                                  ║\n", remaining)
			break
		}
		fmt.Printf("║ %3d. %-54s ║\n", i+1, truncate(chapter.Title, 54))
	}

	fmt.Println("╚════════════════════════════════════════════════════════════╝")

	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s + strings.Repeat(" ", maxLen-len(s))
	}
	return s[:maxLen-3] + "..."
}
