package panel

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var alMuted = lipgloss.NewStyle().Foreground(lipgloss.Color("#6e6a86"))

// Alacritty displays the alacritty config template.
type Alacritty struct {
	width, height int
	content       string
	viewport      viewport.Model
	ready         bool
}

func NewAlacritty() *Alacritty { return &Alacritty{} }

func (a *Alacritty) Title() string     { return "Alacritty" }
func (a *Alacritty) ShortHelp() string { return "j/k: scroll  q: quit  ?: help" }

func (a *Alacritty) Init() tea.Cmd { return nil }

func (a *Alacritty) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SetSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.viewport.Width = msg.Width
		a.viewport.Height = msg.Height
		a.ready = true
		a.updateViewport()
	case AlacrittyContentMsg:
		a.content = msg.Content
		a.updateViewport()
	case tea.KeyMsg:
		var cmd tea.Cmd
		a.viewport, cmd = a.viewport.Update(msg)
		return a, cmd
	}
	return a, nil
}

func (a *Alacritty) updateViewport() {
	if !a.ready {
		return
	}
	if a.content == "" {
		a.viewport.SetContent(alMuted.Render("  Loading..."))
	} else {
		a.viewport.SetContent(a.content)
	}
}

func (a *Alacritty) View() string {
	if !a.ready {
		return alMuted.Render("  Loading...")
	}
	return a.viewport.View()
}
