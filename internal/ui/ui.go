package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/colearendt/pollcat/internal/model"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	okStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	errStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	clockStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("222"))
)

// Model is the Bubble Tea model for the polling TUI.
type Model struct {
	resultsCh    chan model.Result
	results      []model.Result
	lastByTarget map[string]model.Result
	clock        time.Time
	quitting     bool
}

// New creates a new UI Model listening on the provided channel.
func New(resultsCh chan model.Result) Model {
	return Model{
		resultsCh:    resultsCh,
		lastByTarget: make(map[string]model.Result),
		clock:        time.Now(),
	}
}

// Init starts the tick timer and result listener.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tick(),
		waitForResult(m.resultsCh),
	)
}

// Update handles incoming messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		}
	case time.Time:
		m.clock = msg
		return m, tick()
	case model.Result:
		m.results = append(m.results, msg)
		m.lastByTarget[msg.Target] = msg
		return m, waitForResult(m.resultsCh)
	}
	return m, nil
}

// View renders the UI.
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render("pollcat") + "\n")
	b.WriteString(clockStyle.Render(m.clock.Format("15:04:05")) + "\n\n")

	if len(m.lastByTarget) == 0 {
		b.WriteString("Waiting for first poll...\n")
	} else {
		for target, r := range m.lastByTarget {
			status := okStyle.Render("OK")
			if !r.Success {
				status = errStyle.Render("FAIL")
			}
			line := fmt.Sprintf("%-6s %-30s %s %12s  %s",
				r.Type,
				target,
				status,
				r.Latency.Round(time.Microsecond),
				r.Response,
			)
			if r.Error != "" {
				line += fmt.Sprintf("  (%s)", r.Error)
			}
			b.WriteString(line + "\n")
		}
	}

	b.WriteString("\nPress q to quit.\n")
	return b.String()
}

func tick() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return t
	})
}

func waitForResult(ch chan model.Result) tea.Cmd {
	return func() tea.Msg {
		return <-ch
	}
}
