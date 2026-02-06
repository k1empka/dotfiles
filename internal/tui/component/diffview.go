package component

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	diffAdd    = lipgloss.NewStyle().Foreground(lipgloss.Color("#9ccfd8"))
	diffRemove = lipgloss.NewStyle().Foreground(lipgloss.Color("#eb6f92"))
	diffHunk   = lipgloss.NewStyle().Foreground(lipgloss.Color("#c4a7e7"))
	diffMeta   = lipgloss.NewStyle().Foreground(lipgloss.Color("#6e6a86"))
	diffNormal = lipgloss.NewStyle().Foreground(lipgloss.Color("#e0def4"))
)

// DiffView is a viewport that colorizes diff output.
type DiffView struct {
	viewport viewport.Model
	ready    bool
}

// NewDiffView creates a new DiffView.
func NewDiffView() DiffView {
	return DiffView{}
}

// SetSize updates the diff viewport dimensions.
func (d *DiffView) SetSize(w, h int) {
	d.viewport.Width = w
	d.viewport.Height = h
	d.ready = true
}

// SetContent sets and colorizes the diff content.
func (d *DiffView) SetContent(diff string) {
	d.viewport.SetContent(colorizeDiff(diff))
}

func colorizeDiff(diff string) string {
	lines := strings.Split(diff, "\n")
	var b strings.Builder
	for i, line := range lines {
		switch {
		case strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---"):
			b.WriteString(diffMeta.Render(line))
		case strings.HasPrefix(line, "@@"):
			b.WriteString(diffHunk.Render(line))
		case strings.HasPrefix(line, "+"):
			b.WriteString(diffAdd.Render(line))
		case strings.HasPrefix(line, "-"):
			b.WriteString(diffRemove.Render(line))
		default:
			b.WriteString(diffNormal.Render(line))
		}
		if i < len(lines)-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

func (d DiffView) Init() tea.Cmd { return nil }

func (d DiffView) Update(msg tea.Msg) (DiffView, tea.Cmd) {
	var cmd tea.Cmd
	d.viewport, cmd = d.viewport.Update(msg)
	return d, cmd
}

func (d DiffView) View() string {
	if !d.ready {
		return ""
	}
	return d.viewport.View()
}
