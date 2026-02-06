package component

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var errDisplayStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#e0def4")).
	Background(lipgloss.Color("#eb6f92")).
	Padding(0, 1)

type clearErrorMsg struct{}

// ErrorDisplay shows auto-dismissing error messages.
type ErrorDisplay struct {
	message string
	visible bool
}

// NewErrorDisplay creates a new ErrorDisplay.
func NewErrorDisplay() ErrorDisplay {
	return ErrorDisplay{}
}

// Show displays an error message that auto-dismisses after 5 seconds.
func (e *ErrorDisplay) Show(msg string) tea.Cmd {
	e.message = msg
	e.visible = true
	return tea.Tick(5*time.Second, func(time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

// Visible returns whether the error is displayed.
func (e ErrorDisplay) Visible() bool { return e.visible }

func (e ErrorDisplay) Init() tea.Cmd { return nil }

func (e ErrorDisplay) Update(msg tea.Msg) (ErrorDisplay, tea.Cmd) {
	if _, ok := msg.(clearErrorMsg); ok {
		e.visible = false
	}
	return e, nil
}

func (e ErrorDisplay) View() string {
	if !e.visible {
		return ""
	}
	return errDisplayStyle.Render(e.message)
}
