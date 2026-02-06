package component

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	sbStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#1f1d2e")).
		Foreground(lipgloss.Color("#908caa")).
		Padding(0, 1)
	sbKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9ccfd8"))
	sbDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6e6a86"))
)

// StatusBar renders a context-sensitive help/status bar.
type StatusBar struct {
	width int
}

// NewStatusBar creates a new StatusBar.
func NewStatusBar() StatusBar {
	return StatusBar{}
}

// SetWidth updates the status bar width.
func (s *StatusBar) SetWidth(w int) {
	s.width = w
}

// Render renders the status bar with the given help text and optional error.
func (s StatusBar) Render(help string, err error) string {
	errStr := ""
	if err != nil {
		errStr = lipgloss.NewStyle().Foreground(lipgloss.Color("#eb6f92")).Render(" " + err.Error() + " ")
	}
	gap := s.width - lipgloss.Width(help) - lipgloss.Width(errStr)
	if gap < 0 {
		gap = 0
	}
	return sbStyle.Width(s.width).Render(help + strings.Repeat(" ", gap) + errStr)
}

// FormatHelp formats key=description pairs into a help string.
func FormatHelp(pairs ...string) string {
	var parts []string
	for i := 0; i+1 < len(pairs); i += 2 {
		parts = append(parts, sbKeyStyle.Render(pairs[i])+sbDescStyle.Render(" "+pairs[i+1]))
	}
	return strings.Join(parts, "  ")
}
