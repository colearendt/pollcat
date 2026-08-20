package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/colearendt/cli-conn/internal/model"
	"github.com/colearendt/cli-conn/internal/report"
)

var (
	reportFile    string
	reportFormat  string
	reportTargets []string
)

func init() {
	reportCmd.Flags().StringVarP(&reportFile, "file", "f", "", "Input JSON results file (required)")
	reportCmd.Flags().StringVarP(&reportFormat, "format", "t", "table", "Report format: table, csv, json, summary")
	reportCmd.Flags().StringArrayVar(&reportTargets, "target", nil, "Filter results to specific target(s)")
	_ = reportCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(reportCmd)
}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a report from a JSON results file.",
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := os.Open(reportFile)
		if err != nil {
			return fmt.Errorf("open file: %w", err)
		}
		defer f.Close()

		var results []model.Result
		if err := json.NewDecoder(f).Decode(&results); err != nil {
			return fmt.Errorf("decode JSON: %w", err)
		}

		results = report.FilterByTarget(results, reportTargets)

		format := report.Format(reportFormat)
		gen := report.New()
		return gen.WriteReport(os.Stdout, results, format)
	},
}
