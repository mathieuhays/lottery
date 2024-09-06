package start

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	NextView      key.Binding
	PreviousField key.Binding
	NextField     key.Binding
}

var ViewKeyMap = KeyMap{
	NextView: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "start"),
	),
	PreviousField: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "previous field"),
	),
	NextField: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next field"),
	),
}

func ViewKeys() []key.Binding {
	return []key.Binding{
		ViewKeyMap.NextView,
		ViewKeyMap.PreviousField,
		ViewKeyMap.NextField,
	}
}
