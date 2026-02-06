package panel

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	thHeader  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#c4a7e7"))
	thActive  = lipgloss.NewStyle().Foreground(lipgloss.Color("#f6c177"))
	thCursor  = lipgloss.NewStyle().Foreground(lipgloss.Color("#9ccfd8"))
	thNormal  = lipgloss.NewStyle().Foreground(lipgloss.Color("#e0def4"))
	thMuted   = lipgloss.NewStyle().Foreground(lipgloss.Color("#6e6a86"))
	thConfirm = lipgloss.NewStyle().Foreground(lipgloss.Color("#9ccfd8"))
	thDeny    = lipgloss.NewStyle().Foreground(lipgloss.Color("#eb6f92"))
)

// Themes displays available OMP themes with color previews.
type Themes struct {
	width, height int
	themes        []string
	selected      int
	active        string
	colors        map[string][]string
	confirming    bool
}

func NewThemes() *Themes {
	return &Themes{
		colors: make(map[string][]string),
	}
}

func (t *Themes) Title() string { return "Themes" }

func (t *Themes) ShortHelp() string {
	if t.confirming {
		return "y: confirm  n: cancel"
	}
	return "j/k: navigate  enter: apply theme  q: quit  ?: help"
}

func (t *Themes) Init() tea.Cmd { return nil }

func (t *Themes) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SetSizeMsg:
		t.width = msg.Width
		t.height = msg.Height
	case ConfigDataMsg:
		t.active = msg.Theme
		if len(t.themes) == 0 {
			t.themes = []string{"rose-pine", "catppuccin", "tokyo-night"}
		}
	case ThemeColorsMsg:
		t.colors[msg.Theme] = msg.Colors
	case ConfigUpdatedMsg:
		t.confirming = false
	case tea.KeyMsg:
		if t.confirming {
			switch msg.String() {
			case "y", "Y":
				t.confirming = false
				theme := t.themes[t.selected]
				return t, func() tea.Msg { return ThemeApplyMsg{Theme: theme} }
			case "n", "N", "esc":
				t.confirming = false
			}
			return t, nil
		}
		switch msg.String() {
		case "j", "down":
			if t.selected < len(t.themes)-1 {
				t.selected++
			}
		case "k", "up":
			if t.selected > 0 {
				t.selected--
			}
		case "enter":
			if len(t.themes) > 0 {
				t.confirming = true
			}
		}
	}
	return t, nil
}

func (t *Themes) View() string {
	var b strings.Builder

	b.WriteString(thHeader.Render("Available Themes"))
	b.WriteString("\n\n")

	if len(t.themes) == 0 {
		b.WriteString(thMuted.Render("  Loading themes..."))
		return b.String()
	}

	listWidth := 30
	contentWidth := t.width - listWidth - 4
	if contentWidth < 20 {
		contentWidth = 20
	}

	// Theme list.
	var list strings.Builder
	for i, theme := range t.themes {
		cursor := "  "
		if i == t.selected {
			cursor = thCursor.Render("> ")
		}
		name := theme
		style := thNormal
		if i == t.selected {
			style = style.Bold(true)
		}
		suffix := ""
		if theme == t.active {
			suffix = thActive.Render(" (active)")
		}
		list.WriteString(cursor + style.Render(name) + suffix + "\n")
	}

	// Color preview for selected theme.
	var preview strings.Builder
	preview.WriteString(thHeader.Render("Color Palette") + "\n\n")
	if t.selected < len(t.themes) {
		theme := t.themes[t.selected]
		if colors, ok := t.colors[theme]; ok && len(colors) > 0 {
			for _, c := range colors {
				swatch := lipgloss.NewStyle().Background(lipgloss.Color(c)).Render("  ")
				label := thMuted.Render(" " + c)
				preview.WriteString("  " + swatch + label + "\n")
			}
		} else {
			preview.WriteString(thMuted.Render("  Loading colors..."))
		}
	}

	// Confirmation dialog.
	if t.confirming {
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().Width(listWidth).Render(list.String()),
			"  ",
			lipgloss.NewStyle().Width(contentWidth).Render(preview.String()),
		))
		b.WriteString("\n\n")
		b.WriteString("  Apply theme " + thNormal.Bold(true).Render(t.themes[t.selected]) + "? ")
		b.WriteString(thConfirm.Render("[y]") + thMuted.Render(" yes  "))
		b.WriteString(thDeny.Render("[n]") + thMuted.Render(" no"))
		return b.String()
	}

	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(listWidth).Render(list.String()),
		"  ",
		lipgloss.NewStyle().Width(contentWidth).Render(preview.String()),
	))

	return b.String()
}
