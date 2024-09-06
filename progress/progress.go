package progress

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"math/rand"
	"slices"
	"strconv"
	"strings"
	"time"
)

const progressBufferSize = 100

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	labelStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	valueStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("33"))
	logTicketStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("77"))
	logDrawStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	winningStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Blink(true)

	bufferStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(darkGray).PaddingLeft(1)
)

type Model struct {
	buffer        []string
	chances       int
	interval      time.Duration
	cost          float64
	start         time.Time
	end           *time.Time
	totalCost     float64
	ticketCount   int
	winningNumber int
	running       bool
}

type TickMsg struct {
	ticket int
	number int
}

type StartMsg struct {
	chances  int
	interval time.Duration
	cost     float64
}

func New() Model {
	return Model{
		buffer:  make([]string, progressBufferSize),
		running: false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Start(chances int, interval time.Duration, cost float64) tea.Cmd {
	return func() tea.Msg {
		return StartMsg{
			chances:  chances,
			interval: interval,
			cost:     cost,
		}
	}
}

func (m Model) tick() tea.Cmd {
	return tea.Tick(m.interval, func(_ time.Time) tea.Msg {
		return TickMsg{
			ticket: m.ticketCount + 1,
			number: rand.Intn(m.chances),
		}
	})
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case StartMsg:
		m.chances = msg.chances
		m.interval = msg.interval
		m.cost = msg.cost
		m.winningNumber = rand.Intn(m.chances)
		m.running = true
		m.start = time.Now()

		return m, m.tick()
	case TickMsg:
		winning := ""
		if msg.number == m.winningNumber {
			winning = winningStyle.Render(" >>>> Winning!")
			m.running = false

			now := time.Now()
			m.end = &now
		}

		m.ticketCount = msg.ticket
		m.buffer = slices.Concat(m.buffer[1:], []string{
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				logTicketStyle.Render(fmt.Sprintf("Ticket #%s - ", formatTicketNumber(msg.ticket, m.chances))),
				logDrawStyle.Render(formatTicketNumber(msg.number, m.chances)),
				winning,
			),
		})

		if m.running {
			return m, m.tick()
		}

		return m, nil
	}

	return m, nil
}

func prettyCost(cost float64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%.2f", cost)
}

func formatTicketNumber(number int, chances int) string {
	numStr := strconv.Itoa(number)
	chancesStr := strconv.Itoa(chances)
	numLen := len(numStr)
	minLength := len(chancesStr)
	padding := ""

	if numLen < minLength {
		padding = strings.Repeat("0", minLength-numLen)
	}

	return padding + numStr
}

func timeToCoverOdds(chances int, interval time.Duration) time.Duration {
	totalMs := int64(chances) * interval.Milliseconds()
	dur := time.Duration(totalMs) * time.Millisecond

	if dur > time.Hour*24 {
		return dur.Round(time.Hour)
	}

	if dur > time.Hour {
		return dur.Round(time.Minute)
	}

	return dur.Round(time.Second)
}

func rowRenderer(width int) func(label, value string) string {
	rowStyle := lipgloss.NewStyle().Width(width).PaddingRight(1).MarginBottom(1).Align(lipgloss.Left)

	return func(label, value string) string {
		return rowStyle.Render(
			lipgloss.JoinVertical(lipgloss.Left, labelStyle.Render(label+":"), valueStyle.Render(value)),
		)
	}
}

func (m Model) View(width, height int) string {
	if height < 0 {
		return "" // window size not yet set
	}

	var dur time.Duration

	if m.end != nil {
		dur = m.end.Sub(m.start)
	} else {
		dur = time.Since(m.start)
	}

	bufOffset := progressBufferSize - height
	if bufOffset < 0 {
		bufOffset = 0
	}

	infoWidth := width / 3
	row := rowRenderer(infoWidth)

	info := []string{
		row("Chances", fmt.Sprintf("1/%v", m.chances)),
		row("Interval", fmt.Sprintf("%s", m.interval)),
		row("Cost", fmt.Sprintf("%v", m.cost)),
		row("ETA", fmt.Sprintf("%s", timeToCoverOdds(m.chances, m.interval))),
		row("Estimated Cost", prettyCost(m.cost*float64(m.chances))),
		row("Execution Time", fmt.Sprintf("%s", dur.Round(time.Second))),
		row("Total Cost", fmt.Sprintf("%v", prettyCost(m.cost*float64(m.ticketCount)))),
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		lipgloss.JoinVertical(lipgloss.Left, info...),
		bufferStyle.Width(width-infoWidth).Render(lipgloss.JoinVertical(lipgloss.Left, m.buffer[bufOffset:]...)),
	)
}
