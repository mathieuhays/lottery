package main

import (
	"github.com/charmbracelet/bubbles/key"
	"slices"
)

// = KeyMap ===
type defaultKeyMap struct {
	Help          key.Binding
	Quit          key.Binding
	Start         key.Binding
	ExtraKeys     []key.Binding
	ExtraBaseKeys []key.Binding
}

func (k defaultKeyMap) ShortHelp() []key.Binding {
	return slices.Concat(k.ExtraBaseKeys, []key.Binding{k.Help, k.Quit})
}

func (k defaultKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.ExtraKeys,
		k.ShortHelp(),
	}
}

var defaultKeys = defaultKeyMap{
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Start: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "start"),
	),
	ExtraKeys:     make([]key.Binding, 0),
	ExtraBaseKeys: make([]key.Binding, 0),
}
