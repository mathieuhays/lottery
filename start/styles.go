package start

import "github.com/charmbracelet/lipgloss"

var (
	subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "2"}

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 1).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)
)
