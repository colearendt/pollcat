package report

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/colearendt/pollcat/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sampleResults() []model.Result {
	return []model.Result{
		{
			Timestamp: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			Type:      model.PollTypeTCP,
			Target:    "1.2.3.4:80",
			Success:   true,
			Latency:   10 * time.Millisecond,
			Response:  "connected",
		},
		{
			Timestamp: time.Date(2024, 1, 1, 12, 0, 1, 0, time.UTC),
			Type:      model.PollTypeDNS,
			Target:    "example.com",
			Success:   false,
			Latency:   50 * time.Millisecond,
			Error:     "NXDOMAIN",
		},
		{
			Timestamp: time.Date(2024, 1, 1, 12, 0, 2, 0, time.UTC),
			Type:      model.PollTypeTCP,
			Target:    "1.2.3.4:80",
			Success:   true,
			Latency:   20 * time.Millisecond,
			Response:  "connected",
		},
	}
}

func TestGenerator_WriteReport_Table(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	gen := New()
	err := gen.WriteReport(&buf, sampleResults(), FormatTable)
	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "1.2.3.4:80")
	assert.Contains(t, out, "example.com")
	assert.Contains(t, out, "NXDOMAIN")
}

func TestGenerator_WriteReport_CSV(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	gen := New()
	err := gen.WriteReport(&buf, sampleResults(), FormatCSV)
	require.NoError(t, err)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	assert.Len(t, lines, 4) // header + 3 rows
	assert.Equal(t, "timestamp,type,target,success,latency_ms,response,error", lines[0])
	assert.Contains(t, lines[1], "tcp")
	assert.Contains(t, lines[2], "dns")
}

func TestGenerator_WriteReport_JSON(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	gen := New()
	err := gen.WriteReport(&buf, sampleResults(), FormatJSON)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "1.2.3.4:80")
	assert.Contains(t, buf.String(), "example.com")
}

func TestGenerator_WriteReport_Summary(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	gen := New()
	err := gen.WriteReport(&buf, sampleResults(), FormatSummary)
	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "1.2.3.4:80")
	assert.Contains(t, out, "example.com")
	assert.Contains(t, out, "TOTAL")
	assert.Contains(t, out, "AVG")
	assert.Contains(t, out, "NXDOMAIN")
}

func TestGenerator_WriteReport_Unknown(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	gen := New()
	err := gen.WriteReport(&buf, sampleResults(), Format("unknown"))
	assert.EqualError(t, err, "unknown report format: unknown")
}

func TestFilterByTarget(t *testing.T) {
	t.Parallel()
	results := sampleResults()
	filtered := FilterByTarget(results, []string{"1.2.3.4:80"})
	assert.Len(t, filtered, 2)
	for _, r := range filtered {
		assert.Equal(t, "1.2.3.4:80", r.Target)
	}

	filtered = FilterByTarget(results, nil)
	assert.Len(t, filtered, 3)

	filtered = FilterByTarget(results, []string{"missing"})
	assert.Len(t, filtered, 0)
}
