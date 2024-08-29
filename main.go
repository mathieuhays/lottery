package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mathieuhays/lottery/progress"
	"github.com/mathieuhays/lottery/start"
	"log"
	"time"
)

const (
	ViewStart = iota
	ViewProgress
)

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var headerStyle = lipgloss.NewStyle().Foreground(hotPink)

type model struct {
	view     int
	start    start.Model
	progress progress.Model
}

func initialModel() model {
	return model{
		view:     ViewStart,
		start:    start.New(),
		progress: progress.New(),
	}
}

func (m model) Init() tea.Cmd {
	return m.start.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, key.NewBinding(key.WithKeys("q", "ctrl+c"))) {
			return m, tea.Quit
		}
	case start.ViewChangeMsg:
		m.view = ViewProgress
		chances := msg.Chances
		interval := time.Duration(msg.Interval) * time.Millisecond
		cost := msg.Cost

		return m, m.progress.Start(chances, interval, cost)
	}

	var cmd tea.Cmd

	switch m.view {
	case ViewStart:
		m.start, cmd = m.start.Update(msg)
	case ViewProgress:
		m.progress, cmd = m.progress.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	var view string
	header := headerStyle.Width(100).PaddingTop(1).Render("Lottery Simulator")

	switch m.view {
	case ViewStart:
		view = m.start.View()
	case ViewProgress:
		view = m.progress.View()
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, view)
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		log.Fatalf("error starting program: %s", err)
	}
}
