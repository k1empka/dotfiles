package component

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	formLabel    = lipgloss.NewStyle().Foreground(lipgloss.Color("#908caa")).Width(12)
	formActive   = lipgloss.NewStyle().Foreground(lipgloss.Color("#9ccfd8"))
	formInactive = lipgloss.NewStyle().Foreground(lipgloss.Color("#6e6a86"))
)

// FormField describes a single input field.
type FormField struct {
	Label       string
	Placeholder string
	Value       string
}

// FormSubmitMsg is emitted when the form is submitted.
type FormSubmitMsg struct {
	Values map[string]string
}

// FormCancelMsg is emitted when the form is cancelled.
type FormCancelMsg struct{}

// Form is an editable form with labeled text inputs.
type Form struct {
	fields []FormField
	inputs []textinput.Model
	focus  int
}

// NewForm creates a new Form from field definitions.
func NewForm(fields []FormField) Form {
	inputs := make([]textinput.Model, len(fields))
	for i, f := range fields {
		ti := textinput.New()
		ti.Placeholder = f.Placeholder
		ti.SetValue(f.Value)
		ti.CharLimit = 256
		ti.Width = 40
		if i == 0 {
			ti.Focus()
		}
		inputs[i] = ti
	}
	return Form{
		fields: fields,
		inputs: inputs,
	}
}

func (f Form) Init() tea.Cmd {
	return textinput.Blink
}

func (f Form) Update(msg tea.Msg) (Form, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			f.inputs[f.focus].Blur()
			f.focus = (f.focus + 1) % len(f.inputs)
			f.inputs[f.focus].Focus()
			return f, textinput.Blink
		case "shift+tab", "up":
			f.inputs[f.focus].Blur()
			f.focus = (f.focus - 1 + len(f.inputs)) % len(f.inputs)
			f.inputs[f.focus].Focus()
			return f, textinput.Blink
		case "enter":
			values := make(map[string]string)
			for i, field := range f.fields {
				values[field.Label] = f.inputs[i].Value()
			}
			return f, func() tea.Msg { return FormSubmitMsg{Values: values} }
		case "esc":
			return f, func() tea.Msg { return FormCancelMsg{} }
		}
	}

	// Update focused input.
	var cmd tea.Cmd
	f.inputs[f.focus], cmd = f.inputs[f.focus].Update(msg)
	return f, cmd
}

func (f Form) View() string {
	var b strings.Builder
	for i, field := range f.fields {
		label := formLabel.Render(field.Label + ":")
		input := f.inputs[i].View()
		b.WriteString("  " + label + " " + input + "\n")
	}
	b.WriteString("\n")
	b.WriteString(formActive.Render("  enter") + formInactive.Render(" submit  "))
	b.WriteString(formActive.Render("esc") + formInactive.Render(" cancel"))
	return b.String()
}
