package panel

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	ovLabel  = lipgloss.NewStyle().Foreground(lipgloss.Color("#908caa")).Width(14)
	ovValue  = lipgloss.NewStyle().Foreground(lipgloss.Color("#e0def4"))
	ovHeader = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#c4a7e7"))
	ovEdit   = lipgloss.NewStyle().Foreground(lipgloss.Color("#9ccfd8"))
	ovMuted  = lipgloss.NewStyle().Foreground(lipgloss.Color("#6e6a86"))
)

// EditSubmitMsg signals the app to write config and apply.
type EditSubmitMsg struct {
	Values map[string]string
}

// Overview displays configuration summary and supports editing.
type Overview struct {
	width, height int
	name          string
	email         string
	theme         string
	os            string
	version       string
	editing       bool
	inputs        []textinput.Model
	inputLabels   []string
	focusIdx      int
	status        string
}

func NewOverview() *Overview {
	return &Overview{}
}

func (o *Overview) Title() string { return "Overview" }

func (o *Overview) ShortHelp() string {
	if o.editing {
		return "tab: next field  enter: save  esc: cancel"
	}
	return "e: edit  q: quit  ?: help"
}

func (o *Overview) Init() tea.Cmd { return nil }

func (o *Overview) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SetSizeMsg:
		o.width = msg.Width
		o.height = msg.Height
	case ConfigDataMsg:
		o.name = msg.Name
		o.email = msg.Email
		o.theme = msg.Theme
		o.os = msg.OS
		o.version = msg.Version
	case ConfigUpdatedMsg:
		if msg.Err != nil {
			o.status = fmt.Sprintf("Error: %v", msg.Err)
		} else {
			o.status = "Configuration updated successfully"
			o.editing = false
		}
	case tea.KeyMsg:
		if o.editing {
			return o.updateEditing(msg)
		}
		switch msg.String() {
		case "e":
			o.startEditing()
			return o, textinput.Blink
		}
	}

	if o.editing {
		return o.updateInputs(msg)
	}
	return o, nil
}

func (o *Overview) startEditing() {
	o.inputLabels = []string{"name", "email", "theme"}
	values := []string{o.name, o.email, o.theme}
	o.inputs = make([]textinput.Model, len(o.inputLabels))
	for i, v := range values {
		ti := textinput.New()
		ti.SetValue(v)
		ti.CharLimit = 256
		ti.Width = 40
		o.inputs[i] = ti
	}
	o.focusIdx = 0
	o.inputs[0].Focus()
	o.editing = true
	o.status = ""
}

func (o *Overview) updateEditing(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab", "down":
		o.inputs[o.focusIdx].Blur()
		o.focusIdx = (o.focusIdx + 1) % len(o.inputs)
		o.inputs[o.focusIdx].Focus()
		return o, textinput.Blink
	case "shift+tab", "up":
		o.inputs[o.focusIdx].Blur()
		o.focusIdx = (o.focusIdx - 1 + len(o.inputs)) % len(o.inputs)
		o.inputs[o.focusIdx].Focus()
		return o, textinput.Blink
	case "enter":
		values := make(map[string]string)
		for i, label := range o.inputLabels {
			values[label] = o.inputs[i].Value()
		}
		o.editing = false
		return o, func() tea.Msg { return EditSubmitMsg{Values: values} }
	case "esc":
		o.editing = false
		return o, nil
	}
	return o.updateInputs(msg)
}

func (o *Overview) updateInputs(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	o.inputs[o.focusIdx], cmd = o.inputs[o.focusIdx].Update(msg)
	return o, cmd
}

func (o *Overview) View() string {
	var b strings.Builder

	if o.editing {
		b.WriteString(ovHeader.Render("Edit Configuration"))
		b.WriteString("\n\n")
		for i, label := range o.inputLabels {
			b.WriteString("  " + ovLabel.Render(label+":") + " " + o.inputs[i].View() + "\n")
		}
		b.WriteString("\n")
		b.WriteString("  " + ovEdit.Render("enter") + ovMuted.Render(" save  "))
		b.WriteString(ovEdit.Render("esc") + ovMuted.Render(" cancel"))
		return b.String()
	}

	b.WriteString(ovHeader.Render("Configuration"))
	b.WriteString("\n\n")
	rows := []struct{ label, value string }{
		{"Name", o.name},
		{"Email", o.email},
		{"Theme", o.theme},
		{"OS", o.os},
		{"Chezmoi", o.version},
	}
	for _, r := range rows {
		val := r.value
		if val == "" {
			val = ovMuted.Render("(not set)")
		}
		b.WriteString("  " + ovLabel.Render(r.label+":") + " " + ovValue.Render(val) + "\n")
	}

	if o.status != "" {
		b.WriteString("\n  " + ovMuted.Render(o.status))
	}

	return b.String()
}
