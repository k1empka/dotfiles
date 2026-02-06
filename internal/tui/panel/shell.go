package panel

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	shSubActive   = lipgloss.NewStyle().Foreground(lipgloss.Color("#e0def4")).Bold(true).Padding(0, 1).Underline(true)
	shSubInactive = lipgloss.NewStyle().Foreground(lipgloss.Color("#6e6a86")).Padding(0, 1)
	shMuted       = lipgloss.NewStyle().Foreground(lipgloss.Color("#6e6a86"))
)

// Shell displays shell configuration files with sub-tabs.
type Shell struct {
	width, height int
	subTab        int
	subTabs       []string
	contents      [3]string
	viewport      viewport.Model
	ready         bool
}

func NewShell() *Shell {
	return &Shell{
		subTabs: []string{"ZSH", "Bash", "PowerShell"},
	}
}

func (s *Shell) Title() string     { return "Shell" }
func (s *Shell) ShortHelp() string { return "h/l: sub-tab  j/k: scroll  q: quit  ?: help" }

func (s *Shell) Init() tea.Cmd { return nil }

func (s *Shell) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SetSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		s.viewport.Width = msg.Width
		s.viewport.Height = msg.Height - 3 // sub-tab bar + spacing
		s.ready = true
		s.updateViewport()
	case ShellContentMsg:
		if msg.Index >= 0 && msg.Index < 3 {
			s.contents[msg.Index] = msg.Content
		}
		s.updateViewport()
	case tea.KeyMsg:
		switch msg.String() {
		case "h", "left":
			s.subTab = (s.subTab - 1 + len(s.subTabs)) % len(s.subTabs)
			s.updateViewport()
			return s, nil
		case "l", "right":
			s.subTab = (s.subTab + 1) % len(s.subTabs)
			s.updateViewport()
			return s, nil
		}
		var cmd tea.Cmd
		s.viewport, cmd = s.viewport.Update(msg)
		return s, cmd
	}
	return s, nil
}

func (s *Shell) updateViewport() {
	if !s.ready {
		return
	}
	content := s.contents[s.subTab]
	if content == "" {
		content = shMuted.Render("  Loading...")
	}
	s.viewport.SetContent(content)
}

func (s *Shell) View() string {
	var b strings.Builder

	// Sub-tab bar.
	for i, t := range s.subTabs {
		if i == s.subTab {
			b.WriteString(shSubActive.Render(t))
		} else {
			b.WriteString(shSubInactive.Render(t))
		}
	}
	b.WriteString("\n\n")

	if !s.ready {
		b.WriteString(shMuted.Render("  Loading..."))
		return b.String()
	}

	b.WriteString(s.viewport.View())
	return b.String()
}
