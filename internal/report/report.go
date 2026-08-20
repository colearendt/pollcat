package report

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/colearendt/pollcat/internal/model"
)

// Format defines the report output format.
type Format string

const (
	FormatTable   Format = "table"
	FormatCSV     Format = "csv"
	FormatJSON    Format = "json"
	FormatSummary Format = "summary"
)

// Generator writes reports from a slice of results.
type Generator struct{}

// New creates a new Generator.
func New() *Generator {
	return &Generator{}
}

// WriteReport emits a report to w in the requested format.
func (g *Generator) WriteReport(w io.Writer, results []model.Result, format Format) error {
	switch format {
	case FormatTable:
		return g.writeTable(w, results)
	case FormatCSV:
		return g.writeCSV(w, results)
	case FormatJSON:
		return g.writeJSON(w, results)
	case FormatSummary:
		return g.writeSummary(w, results)
	default:
		return fmt.Errorf("unknown report format: %s", format)
	}
}

func (g *Generator) writeTable(w io.Writer, results []model.Result) error {
	fmt.Fprintf(w, "%-24s %-6s %-30s %-8s %-12s %s\n", "TIME", "TYPE", "TARGET", "SUCCESS", "LATENCY", "RESPONSE/ERROR")
	fmt.Fprintln(w, strings.Repeat("-", 100))
	for _, r := range results {
		success := "OK"
		if !r.Success {
			success = "FAIL"
		}
		resp := r.Response
		if resp == "" {
			resp = r.Error
		}
		fmt.Fprintf(w, "%-24s %-6s %-30s %-8s %-12s %s\n",
			r.Timestamp.Format(time.RFC3339), r.Type, r.Target, success, r.Latency, resp)
	}
	return nil
}

func (g *Generator) writeCSV(w io.Writer, results []model.Result) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	if err := cw.Write([]string{"timestamp", "type", "target", "success", "latency_ms", "response", "error"}); err != nil {
		return err
	}
	for _, r := range results {
		latencyMs := float64(r.Latency) / float64(time.Millisecond)
		row := []string{
			r.Timestamp.Format(time.RFC3339Nano),
			string(r.Type),
			r.Target,
			fmt.Sprintf("%t", r.Success),
			fmt.Sprintf("%.3f", latencyMs),
			r.Response,
			r.Error,
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) writeJSON(w io.Writer, results []model.Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}

type targetSummary struct {
	Type       model.PollType
	Target     string
	Total      int
	Successes  int
	Failures   int
	MinLatency time.Duration
	MaxLatency time.Duration
	AvgLatency time.Duration
	LastResp   string
	LastErr    string
}

func (g *Generator) writeSummary(w io.Writer, results []model.Result) error {
	// Group by type+target
	sums := make(map[string]*targetSummary)
	for _, r := range results {
		key := string(r.Type) + ":" + r.Target
		s, ok := sums[key]
		if !ok {
			s = &targetSummary{Type: r.Type, Target: r.Target, MinLatency: r.Latency}
			sums[key] = s
		}
		s.Total++
		if r.Success {
			s.Successes++
		} else {
			s.Failures++
		}
		if r.Latency < s.MinLatency {
			s.MinLatency = r.Latency
		}
		if r.Latency > s.MaxLatency {
			s.MaxLatency = r.Latency
		}
		// accumulate for avg
		s.AvgLatency += r.Latency
		s.LastResp = r.Response
		s.LastErr = r.Error
	}

	for _, s := range sums {
		if s.Total > 0 {
			s.AvgLatency = s.AvgLatency / time.Duration(s.Total)
		}
	}

	// Sort for deterministic output
	keys := make([]string, 0, len(sums))
	for k := range sums {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintf(w, "%-6s %-30s %6s %6s %6s %12s %12s %12s %s\n",
		"TYPE", "TARGET", "TOTAL", "OK", "FAIL", "MIN", "MAX", "AVG", "LAST")
	fmt.Fprintln(w, strings.Repeat("-", 110))
	for _, k := range keys {
		s := sums[k]
		last := s.LastResp
		if last == "" {
			last = s.LastErr
		}
		fmt.Fprintf(w, "%-6s %-30s %6d %6d %6d %12s %12s %12s %s\n",
			s.Type, s.Target, s.Total, s.Successes, s.Failures,
			s.MinLatency.Round(time.Microsecond),
			s.MaxLatency.Round(time.Microsecond),
			s.AvgLatency.Round(time.Microsecond),
			last)
	}
	return nil
}

// FilterByTarget returns only results whose Target field matches one of the given targets.
func FilterByTarget(results []model.Result, targets []string) []model.Result {
	if len(targets) == 0 {
		return results
	}
	allowed := make(map[string]struct{}, len(targets))
	for _, t := range targets {
		allowed[t] = struct{}{}
	}
	var out []model.Result
	for _, r := range results {
		if _, ok := allowed[r.Target]; ok {
			out = append(out, r)
		}
	}
	return out
}
