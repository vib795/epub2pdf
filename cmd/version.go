package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version   = "1.0.0"
	BuildDate = "unknown"
	GitCommit = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("epub2pdf %s\n", Version)
		fmt.Printf("  Build Date: %s\n", BuildDate)
		fmt.Printf("  Git Commit: %s\n", GitCommit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
