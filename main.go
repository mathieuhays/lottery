package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mathieuhays/lottery/header"
	"github.com/mathieuhays/lottery/progress"
	"github.com/mathieuhays/lottery/start"
	"log"
	"time"
)

type model struct {
	width  int
	height int

	keys     defaultKeyMap
	header   header.Model
	help     help.Model
	start    start.Model
	progress progress.Model
	quitting bool

	debug          string
	isProgressView bool
}

func initialModel() model {
	m := model{
		keys:     defaultKeys,
		header:   header.New("Lottery Odd Visualizer"),
		help:     help.New(),
		start:    start.New(),
		progress: progress.New(),
	}

	m.keys.ExtraBaseKeys = start.ViewKeys()

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var mainCmd, startCmd, progressCmd tea.Cmd

	switch msg := message.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Start):
			mainCmd = m.progress.Start(2000, time.Millisecond*100, 1)
			m.debug = "start triggered"
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			mainCmd = tea.Quit
		}

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

	case start.ViewChangeMsg:
		m.isProgressView = true
		m.keys.ExtraBaseKeys = progress.ViewKeys()
		chances := msg.Chances
		interval := time.Duration(msg.Interval) * time.Millisecond
		cost := msg.Cost

		return m, m.progress.Start(chances, interval, cost)
	}

	if m.isProgressView {
		m.progress, progressCmd = m.progress.Update(message)
	} else {
		m.start, startCmd = m.start.Update(message)
	}

	return m, tea.Batch(mainCmd, startCmd, progressCmd)
}

func (m model) View() string {
	if m.height == 0 {
		return "Loading..."
	}

	var content string
	helpStyle := lipgloss.NewStyle().
		Width(m.width).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}).
		PaddingTop(1)

	headerView := m.header.View(m.width)
	helpView := helpStyle.Render(m.help.View(m.keys))
	contentHeight := m.height - lipgloss.Height(headerView) - lipgloss.Height(helpView)

	if m.isProgressView {
		content = m.progress.View(m.width, contentHeight)
	} else {
		content = m.start.View(m.width, contentHeight)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerView,
		content,
		helpView,
	)
}

func main() {
	if _, err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Run(); err != nil {
		log.Fatal(err)
	}
}
