package header

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const hotPink = lipgloss.Color("#FF06B7")

var (
	subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}

	headerStyle = lipgloss.NewStyle().
			Foreground(hotPink).
			PaddingTop(1).
			PaddingBottom(1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(subtle).
			BorderBottom(true)
)

type Model string

func New(title string) Model {
	return Model(title)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(message tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View(width int) string {
	return headerStyle.Width(width).Align(lipgloss.Center).Render(string(m))
}
