package component

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	lineNumStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#6e6a86")).Width(4).Align(lipgloss.Right)
	lineStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#e0def4"))
)

// Viewer is a scrollable text viewport with optional line numbers.
type Viewer struct {
	viewport    viewport.Model
	content     string
	lineNumbers bool
	ready       bool
}

// NewViewer creates a new Viewer.
func NewViewer(lineNumbers bool) Viewer {
	return Viewer{lineNumbers: lineNumbers}
}

// SetSize updates the viewport dimensions.
func (v *Viewer) SetSize(w, h int) {
	v.viewport.Width = w
	v.viewport.Height = h
	v.ready = true
	v.refreshContent()
}

// SetContent updates the displayed text.
func (v *Viewer) SetContent(content string) {
	v.content = content
	v.refreshContent()
}

func (v *Viewer) refreshContent() {
	if !v.ready {
		return
	}
	if v.lineNumbers {
		lines := strings.Split(v.content, "\n")
		var b strings.Builder
		for i, line := range lines {
			b.WriteString(lineNumStyle.Render(fmt.Sprintf("%d", i+1)))
			b.WriteString("  ")
			b.WriteString(lineStyle.Render(line))
			if i < len(lines)-1 {
				b.WriteString("\n")
			}
		}
		v.viewport.SetContent(b.String())
	} else {
		v.viewport.SetContent(v.content)
	}
}

func (v Viewer) Init() tea.Cmd { return nil }

func (v Viewer) Update(msg tea.Msg) (Viewer, tea.Cmd) {
	var cmd tea.Cmd
	v.viewport, cmd = v.viewport.Update(msg)
	return v, cmd
}

func (v Viewer) View() string {
	if !v.ready {
		return ""
	}
	return v.viewport.View()
}

// ScrollPercent returns the current scroll position as a percentage.
func (v Viewer) ScrollPercent() float64 {
	return v.viewport.ScrollPercent()
}
