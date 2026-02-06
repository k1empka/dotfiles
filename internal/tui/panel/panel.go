package panel

import tea "github.com/charmbracelet/bubbletea"

// Panel is the interface that all tab panels implement.
type Panel interface {
	tea.Model
	Title() string
	ShortHelp() string
}
