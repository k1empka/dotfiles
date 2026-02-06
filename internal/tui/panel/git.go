package panel

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var gitMuted = lipgloss.NewStyle().Foreground(lipgloss.Color("#6e6a86"))

// Git displays the gitconfig with resolved template values.
type Git struct {
	width, height int
	content       string
	viewport      viewport.Model
	ready         bool
}

func NewGit() *Git { return &Git{} }

func (g *Git) Title() string     { return "Git" }
func (g *Git) ShortHelp() string { return "j/k: scroll  q: quit  ?: help" }

func (g *Git) Init() tea.Cmd { return nil }

func (g *Git) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SetSizeMsg:
		g.width = msg.Width
		g.height = msg.Height
		g.viewport.Width = msg.Width
		g.viewport.Height = msg.Height
		g.ready = true
		g.updateViewport()
	case GitContentMsg:
		g.content = msg.Content
		g.updateViewport()
	case tea.KeyMsg:
		var cmd tea.Cmd
		g.viewport, cmd = g.viewport.Update(msg)
		return g, cmd
	}
	return g, nil
}

func (g *Git) updateViewport() {
	if !g.ready {
		return
	}
	if g.content == "" {
		g.viewport.SetContent(gitMuted.Render("  Loading..."))
	} else {
		g.viewport.SetContent(g.content)
	}
}

func (g *Git) View() string {
	if !g.ready {
		return gitMuted.Render("  Loading...")
	}
	return g.viewport.View()
}
