package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cli-conn",
	Short: "A lightweight CLI for polling network connectivity and DNS resolution.",
	Long: `cli-conn runs multiple concurrent pollers (TCP connect, DNS lookup)
and displays a real-time TUI with a system clock. After polling, you can
generate CLI-friendly reports or CSV exports.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
