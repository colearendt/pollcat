package ui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/colearendt/cli-conn/internal/model"
)

func TestModel_Init(t *testing.T) {
	ch := make(chan model.Result)
	m := New(ch)
	cmd := m.Init()
	assert.NotNil(t, cmd)
}

func TestModel_Update_Result(t *testing.T) {
	ch := make(chan model.Result)
	m := New(ch)
	res := model.Result{Target: "example.com", Type: model.PollTypeDNS, Success: true, Latency: 10 * time.Millisecond}
	newModel, cmd := m.Update(res)
	assert.NotNil(t, cmd)
	assert.Len(t, newModel.(Model).results, 1)
	assert.Contains(t, newModel.(Model).lastByTarget, "example.com")
}

func TestModel_Update_Tick(t *testing.T) {
	ch := make(chan model.Result)
	m := New(ch)
	now := time.Now()
	newModel, cmd := m.Update(now)
	assert.NotNil(t, cmd)
	assert.Equal(t, now, newModel.(Model).clock)
}

func TestModel_Update_Quit(t *testing.T) {
	ch := make(chan model.Result)
	m := New(ch)
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	newModel, cmd := m.Update(msg)
	assert.NotNil(t, cmd)
	assert.True(t, newModel.(Model).quitting)
}

func TestModel_View(t *testing.T) {
	ch := make(chan model.Result)
	m := New(ch)
	view := m.View()
	assert.Contains(t, view, "cli-conn")
	assert.Contains(t, view, "Waiting for first poll")
}
