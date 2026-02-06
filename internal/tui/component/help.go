package component

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	helpTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#c4a7e7")).
			MarginBottom(1)
	helpKey = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9ccfd8")).
		Width(16)
	helpDesc = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#e0def4"))
	helpSection = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f6c177")).
			Bold(true).
			MarginTop(1)
)

// HelpEntry is a single key binding help entry.
type HelpEntry struct {
	Key  string
	Desc string
}

// HelpSection groups help entries under a heading.
type HelpSection struct {
	Title   string
	Entries []HelpEntry
}

// Help is a full-screen help overlay.
type Help struct {
	visible  bool
	sections []HelpSection
	width    int
	height   int
}

// NewHelp creates a new Help overlay with default sections.
func NewHelp() Help {
	return Help{
		sections: []HelpSection{
			{
				Title: "Global",
				Entries: []HelpEntry{
					{"Tab / Shift+Tab", "Next / previous tab"},
					{"1-6", "Jump to tab"},
					{"q / Ctrl+C", "Quit"},
					{"?", "Toggle help"},
				},
			},
			{
				Title: "Navigation",
				Entries: []HelpEntry{
					{"j / k", "Down / up"},
					{"g / G", "Top / bottom"},
					{"Ctrl+D / Ctrl+U", "Half page down / up"},
					{"h / l", "Left / right (sub-tabs, panes)"},
				},
			},
			{
				Title: "Overview",
				Entries: []HelpEntry{
					{"e", "Edit configuration"},
					{"Enter", "Submit form"},
					{"Esc", "Cancel edit"},
				},
			},
			{
				Title: "Themes",
				Entries: []HelpEntry{
					{"Enter", "Apply selected theme"},
				},
			},
			{
				Title: "Chezmoi",
				Entries: []HelpEntry{
					{"s", "Run status"},
					{"d", "Run diff"},
					{"a", "Run apply (with confirmation)"},
					{"y / n", "Confirm / cancel"},
				},
			},
		},
	}
}

// Toggle shows or hides the help overlay.
func (h *Help) Toggle() {
	h.visible = !h.visible
}

// Visible returns whether help is shown.
func (h Help) Visible() bool { return h.visible }

// SetSize updates the help overlay dimensions.
func (h *Help) SetSize(w, ht int) {
	h.width = w
	h.height = ht
}

func (h Help) Init() tea.Cmd { return nil }

func (h Help) Update(msg tea.Msg) (Help, tea.Cmd) {
	if km, ok := msg.(tea.KeyMsg); ok {
		if km.String() == "?" || km.String() == "esc" || km.String() == "q" {
			h.visible = false
		}
	}
	return h, nil
}

func (h Help) View() string {
	if !h.visible {
		return ""
	}

	var b strings.Builder
	b.WriteString(helpTitle.Render("Keyboard Shortcuts"))
	b.WriteString("\n")

	for _, section := range h.sections {
		b.WriteString(helpSection.Render(section.Title))
		b.WriteString("\n")
		for _, entry := range section.Entries {
			b.WriteString("  " + helpKey.Render(entry.Key) + helpDesc.Render(entry.Desc) + "\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#6e6a86")).Render("Press ? or Esc to close"))

	content := b.String()
	overlay := lipgloss.NewStyle().
		Width(h.width).
		Height(h.height).
		Background(lipgloss.Color("#191724")).
		Padding(1, 2)

	return overlay.Render(content)
}
