package tui

import "github.com/charmbracelet/bubbles/key"

// globalKeyMap defines keys handled by the top-level App.
type globalKeyMap struct {
	Quit     key.Binding
	NextTab  key.Binding
	PrevTab  key.Binding
	Tab1     key.Binding
	Tab2     key.Binding
	Tab3     key.Binding
	Tab4     key.Binding
	Tab5     key.Binding
	Tab6     key.Binding
	Tab7     key.Binding
	Help     key.Binding
	Up       key.Binding
	Down     key.Binding
	Top      key.Binding
	Bottom   key.Binding
	HalfUp   key.Binding
	HalfDown key.Binding
	Left     key.Binding
	Right    key.Binding
	Enter    key.Binding
	Escape   key.Binding
}

var globalKeys = globalKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	NextTab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next tab"),
	),
	PrevTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "prev tab"),
	),
	Tab1: key.NewBinding(key.WithKeys("1")),
	Tab2: key.NewBinding(key.WithKeys("2")),
	Tab3: key.NewBinding(key.WithKeys("3")),
	Tab4: key.NewBinding(key.WithKeys("4")),
	Tab5: key.NewBinding(key.WithKeys("5")),
	Tab6: key.NewBinding(key.WithKeys("6")),
	Tab7: key.NewBinding(key.WithKeys("7")),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("j/k", "up/down"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
	),
	Top: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g/G", "top/bottom"),
	),
	Bottom: key.NewBinding(
		key.WithKeys("G"),
	),
	HalfUp: key.NewBinding(
		key.WithKeys("ctrl+u"),
		key.WithHelp("ctrl+u/d", "half page"),
	),
	HalfDown: key.NewBinding(
		key.WithKeys("ctrl+d"),
	),
	Left: key.NewBinding(
		key.WithKeys("h", "left"),
		key.WithHelp("h/l", "left/right"),
	),
	Right: key.NewBinding(
		key.WithKeys("l", "right"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
}
