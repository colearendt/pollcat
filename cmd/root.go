package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pollcat",
	Short: "A lightweight CLI for polling network connectivity and DNS resolution.",
	Long: `pollcat runs multiple concurrent pollers (TCP connect, DNS lookup)
and displays a real-time TUI with a system clock. After polling, you can
generate CLI-friendly reports or CSV exports.`,
}

// SetVersion injects build-time version info into the root command.
func SetVersion(version, commit, date string) {
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
