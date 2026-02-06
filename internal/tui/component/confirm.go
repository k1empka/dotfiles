package component

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	dialogStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#c4a7e7")).
			Padding(1, 2).
			Width(50)
	dialogTitle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#e0def4"))
	dialogHelp    = lipgloss.NewStyle().Foreground(lipgloss.Color("#6e6a86"))
	dialogConfirm = lipgloss.NewStyle().Foreground(lipgloss.Color("#9ccfd8"))
	dialogDeny    = lipgloss.NewStyle().Foreground(lipgloss.Color("#eb6f92"))
)

// ConfirmMsg is emitted with the user's confirmation choice.
type ConfirmMsg struct {
	Confirmed bool
	Tag       string
}

// Confirm is a y/n confirmation dialog.
type Confirm struct {
	message string
	tag     string
	visible bool
}

// NewConfirm creates a new confirmation dialog.
func NewConfirm() Confirm {
	return Confirm{}
}

// Show displays the dialog with a message and tag for identification.
func (c *Confirm) Show(message, tag string) {
	c.message = message
	c.tag = tag
	c.visible = true
}

// Hide dismisses the dialog.
func (c *Confirm) Hide() {
	c.visible = false
}

// Visible returns whether the dialog is shown.
func (c Confirm) Visible() bool { return c.visible }

func (c Confirm) Init() tea.Cmd { return nil }

func (c Confirm) Update(msg tea.Msg) (Confirm, tea.Cmd) {
	if !c.visible {
		return c, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			c.visible = false
			tag := c.tag
			return c, func() tea.Msg { return ConfirmMsg{Confirmed: true, Tag: tag} }
		case "n", "N", "esc":
			c.visible = false
			tag := c.tag
			return c, func() tea.Msg { return ConfirmMsg{Confirmed: false, Tag: tag} }
		}
	}
	return c, nil
}

func (c Confirm) View() string {
	if !c.visible {
		return ""
	}
	content := dialogTitle.Render(c.message) + "\n\n"
	content += dialogConfirm.Render("[y]") + dialogHelp.Render(" yes  ")
	content += dialogDeny.Render("[n]") + dialogHelp.Render(" no")
	return dialogStyle.Render(content)
}
