package start

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strconv"
	"strings"
)

type errMsg error
type ViewChangeMsg struct {
	Chances  int
	Interval int
	Cost     float64
}

const (
	chances = iota
	interval
	cost
)

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(hotPink)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
)

type Model struct {
	inputs  []textinput.Model
	focused int
	err     error
}

func chancesValidator(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("chances number is required")
	}

	c := strings.ReplaceAll(s, "_", "")
	_, err := strconv.ParseInt(c, 10, 64)

	return err
}

func intervalValidator(s string) error {
	if len(s) == 0 {
		return nil // nothing to validate
	}

	_, err := strconv.Atoi(s)
	return err
}

func costValidator(s string) error {
	if len(s) == 0 {
		return nil // nothing to validate
	}

	_, err := strconv.ParseFloat(s, 64)

	return err
}

func New() Model {
	var inputs []textinput.Model = make([]textinput.Model, 3)

	// chances
	inputs[chances] = textinput.New()
	inputs[chances].Placeholder = "25_000_000"
	inputs[chances].Focus()
	inputs[chances].Prompt = ""
	inputs[chances].Validate = chancesValidator

	// interval
	inputs[interval] = textinput.New()
	inputs[interval].Placeholder = "1000"
	inputs[interval].Prompt = ""
	inputs[interval].Validate = intervalValidator

	// Cost
	inputs[cost] = textinput.New()
	inputs[cost].Placeholder = "0.5"
	inputs[cost].Prompt = ""
	inputs[cost].Validate = costValidator

	return Model{
		inputs:  inputs,
		focused: chances,
		err:     nil,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				return m, m.nextView()
			}

			// next input
			m.focused = nextInput(m)
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.focused = prevInput(m)
		case tea.KeyTab, tea.KeyCtrlN:
			m.focused = nextInput(m)
		}

		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()
	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	errMessage := ""

	if m.err != nil {
		errMessage = "\n\n" + m.err.Error() + "\n\n"
	}

	return fmt.Sprintf(`%s

%s
%s

%s
%s

%s
%s

%s
`,
		errMessage,
		inputStyle.Width(30).Render("Chances"),
		m.inputs[chances].View(),
		inputStyle.Width(30).Render("Interval"),
		m.inputs[interval].View(),
		inputStyle.Width(30).Render("Cost"),
		m.inputs[cost].View(),
		continueStyle.Render("Continue ->"),
	) + "\n"
}

func nextInput(m Model) int {
	return (m.focused + 1) % len(m.inputs)
}

func prevInput(m Model) int {
	f := m.focused - 1

	if f < 0 {
		f = len(m.inputs) - 1
	}

	return f
}

func (m Model) nextView() tea.Cmd {
	chancesStr := strings.ReplaceAll(m.inputs[chances].Value(), "_", "")
	chancesStr = strings.TrimSpace(chancesStr)
	chancesValue, err := strconv.Atoi(chancesStr)
	if err != nil {
		return func() tea.Msg {
			return errMsg(err)
		}
	}

	intervalValue, err := strconv.Atoi(m.inputs[interval].Value())
	if err != nil {
		return func() tea.Msg {
			return errMsg(err)
		}
	}

	costValue, err := strconv.ParseFloat(m.inputs[cost].Value(), 64)
	if err != nil {
		return func() tea.Msg {
			return errMsg(err)
		}
	}

	return func() tea.Msg {
		return ViewChangeMsg{
			Chances:  chancesValue,
			Interval: intervalValue,
			Cost:     costValue,
		}
	}
}
