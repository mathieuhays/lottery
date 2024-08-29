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

const progressBufferSize = 10

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
			Border(lipgloss.NormalBorder(), false, false, false, true).MarginLeft(2).PaddingLeft(2)
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

type tickMsg struct {
	ticket int
	number int
}

type startMsg struct {
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
		return startMsg{
			chances:  chances,
			interval: interval,
			cost:     cost,
		}
	}
}

func (m Model) tick() tea.Cmd {
	return tea.Tick(m.interval, func(_ time.Time) tea.Msg {
		return tickMsg{
			ticket: m.ticketCount + 1,
			number: rand.Intn(m.chances),
		}
	})
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case startMsg:
		m.chances = msg.chances
		m.interval = msg.interval
		m.cost = msg.cost
		m.winningNumber = rand.Intn(m.chances)
		m.running = true
		m.start = time.Now()

		return m, m.tick()
	case tickMsg:
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

func (m Model) View() string {
	var dur time.Duration

	if m.end != nil {
		dur = m.end.Sub(m.start)
	} else {
		dur = time.Since(m.start)
	}

	info := []string{
		labelStyle.Render("Chances: ") + valueStyle.Render(fmt.Sprintf("1/%v", m.chances)),
		labelStyle.Render("Interval: ") + valueStyle.Render(fmt.Sprintf("%s", m.interval)),
		labelStyle.Render("Cost: ") + valueStyle.Render(fmt.Sprintf("%v", m.cost)),
		labelStyle.Render("Time to match odds: ") + valueStyle.Render(fmt.Sprintf("%s", timeToCoverOdds(m.chances, m.interval))),
		labelStyle.Render("Cost to match odds: ") + valueStyle.Render(prettyCost(m.cost*float64(m.chances))),
		"--------------------",
		labelStyle.Render("Execution time: ") + valueStyle.Render(fmt.Sprintf("%s", dur.Round(time.Second))),
		labelStyle.Render("Total cost: ") + valueStyle.Render(fmt.Sprintf("%v", prettyCost(m.cost*float64(m.ticketCount)))),
	}

	return "\n" + lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.JoinVertical(lipgloss.Left, info...),
		bufferStyle.Render(lipgloss.JoinVertical(lipgloss.Left, m.buffer...)),
	)
}
