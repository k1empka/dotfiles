package panel

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	czHeader  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#c4a7e7"))
	czKey     = lipgloss.NewStyle().Foreground(lipgloss.Color("#f6c177"))
	czActive  = lipgloss.NewStyle().Foreground(lipgloss.Color("#9ccfd8")).Bold(true)
	czNormal  = lipgloss.NewStyle().Foreground(lipgloss.Color("#6e6a86"))
	czMuted   = lipgloss.NewStyle().Foreground(lipgloss.Color("#6e6a86"))
	czAdd     = lipgloss.NewStyle().Foreground(lipgloss.Color("#9ccfd8"))
	czRemove  = lipgloss.NewStyle().Foreground(lipgloss.Color("#eb6f92"))
	czHunk    = lipgloss.NewStyle().Foreground(lipgloss.Color("#c4a7e7"))
	czMeta    = lipgloss.NewStyle().Foreground(lipgloss.Color("#6e6a86"))
	czConfirm = lipgloss.NewStyle().Foreground(lipgloss.Color("#9ccfd8"))
	czDeny    = lipgloss.NewStyle().Foreground(lipgloss.Color("#eb6f92"))
)

// Chezmoi provides status, diff, and apply operations.
type Chezmoi struct {
	width, height int
	viewport      viewport.Model
	spinner       spinner.Model
	mode          string
	loading       bool
	confirming    bool
	output        string
	err           error
	ready         bool
}

func NewChezmoi() *Chezmoi {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#c4a7e7"))
	return &Chezmoi{
		mode:    "status",
		spinner: sp,
	}
}

func (c *Chezmoi) Title() string { return "Chezmoi" }

func (c *Chezmoi) ShortHelp() string {
	if c.confirming {
		return "y: confirm apply  n: cancel"
	}
	return "s: status  d: diff  a: apply  j/k: scroll  q: quit  ?: help"
}

func (c *Chezmoi) Init() tea.Cmd { return nil }

func (c *Chezmoi) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SetSizeMsg:
		c.width = msg.Width
		c.height = msg.Height
		c.viewport.Width = msg.Width
		c.viewport.Height = msg.Height - 5
		c.ready = true
		if c.output != "" {
			c.setViewportContent()
		}
	case ChezmoiOutputMsg:
		c.loading = false
		c.mode = msg.Mode
		if msg.Err != nil {
			c.err = msg.Err
			c.output = fmt.Sprintf("Error: %v", msg.Err)
		} else {
			c.err = nil
			c.output = msg.Output
			if c.output == "" {
				c.output = "(no output)"
			}
		}
		c.setViewportContent()
	case spinner.TickMsg:
		if c.loading {
			var cmd tea.Cmd
			c.spinner, cmd = c.spinner.Update(msg)
			return c, cmd
		}
	case tea.KeyMsg:
		if c.confirming {
			switch msg.String() {
			case "y", "Y":
				c.confirming = false
				c.loading = true
				return c, tea.Batch(
					c.spinner.Tick,
					func() tea.Msg { return RunChezmoiMsg{Mode: "apply"} },
				)
			case "n", "N", "esc":
				c.confirming = false
			}
			return c, nil
		}
		switch msg.String() {
		case "s":
			c.mode = "status"
			c.loading = true
			return c, tea.Batch(
				c.spinner.Tick,
				func() tea.Msg { return RunChezmoiMsg{Mode: "status"} },
			)
		case "d":
			c.mode = "diff"
			c.loading = true
			return c, tea.Batch(
				c.spinner.Tick,
				func() tea.Msg { return RunChezmoiMsg{Mode: "diff"} },
			)
		case "a":
			c.confirming = true
			return c, nil
		}
		// Scroll viewport.
		var cmd tea.Cmd
		c.viewport, cmd = c.viewport.Update(msg)
		return c, cmd
	}
	return c, nil
}

func (c *Chezmoi) setViewportContent() {
	if !c.ready {
		return
	}
	if c.mode == "diff" {
		c.viewport.SetContent(colorizeDiff(c.output))
	} else {
		c.viewport.SetContent(c.output)
	}
}

func colorizeDiff(diff string) string {
	lines := strings.Split(diff, "\n")
	var b strings.Builder
	for i, line := range lines {
		switch {
		case strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---"):
			b.WriteString(czMeta.Render(line))
		case strings.HasPrefix(line, "@@"):
			b.WriteString(czHunk.Render(line))
		case strings.HasPrefix(line, "+"):
			b.WriteString(czAdd.Render(line))
		case strings.HasPrefix(line, "-"):
			b.WriteString(czRemove.Render(line))
		default:
			b.WriteString(line)
		}
		if i < len(lines)-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

func (c *Chezmoi) View() string {
	var b strings.Builder

	b.WriteString(czHeader.Render("Chezmoi Operations"))
	b.WriteString("\n\n")

	// Action bar.
	for _, action := range []struct{ key, label string }{
		{"s", "status"},
		{"d", "diff"},
		{"a", "apply"},
	} {
		style := czNormal
		if action.label == c.mode {
			style = czActive
		}
		b.WriteString("  " + czKey.Render("["+action.key+"]") + " " + style.Render(action.label))
	}
	b.WriteString("\n\n")

	if c.confirming {
		b.WriteString("  Apply changes? ")
		b.WriteString(czConfirm.Render("[y]") + czMuted.Render(" yes  "))
		b.WriteString(czDeny.Render("[n]") + czMuted.Render(" no"))
		return b.String()
	}

	if c.loading {
		b.WriteString("  " + c.spinner.View() + " Running...")
		return b.String()
	}

	if c.ready {
		b.WriteString(c.viewport.View())
	}

	return b.String()
}
