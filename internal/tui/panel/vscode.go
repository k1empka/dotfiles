package panel

import (
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	vsHeader = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#c4a7e7"))
	vsCursor = lipgloss.NewStyle().Foreground(lipgloss.Color("#9ccfd8"))
	vsActive = lipgloss.NewStyle().Foreground(lipgloss.Color("#e0def4")).Bold(true)
	vsNormal = lipgloss.NewStyle().Foreground(lipgloss.Color("#908caa"))
	vsMuted  = lipgloss.NewStyle().Foreground(lipgloss.Color("#6e6a86"))
)

// VSCode displays VS Code config files in a split-pane browser.
type VSCode struct {
	width, height int
	files         []string
	selected      int
	content       string
	contentPath   string
	viewport      viewport.Model
	ready         bool
	focusContent  bool
	listOffset    int
}

func NewVSCode() *VSCode { return &VSCode{} }

func (v *VSCode) Title() string { return "VS Code" }

func (v *VSCode) ShortHelp() string {
	if v.focusContent {
		return "h: file list  j/k: scroll  q: quit  ?: help"
	}
	return "j/k: navigate  l/enter: view file  q: quit  ?: help"
}

func (v *VSCode) Init() tea.Cmd { return nil }

func (v *VSCode) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SetSizeMsg:
		v.width = msg.Width
		v.height = msg.Height
		v.viewport.Width = msg.Width - 34
		v.viewport.Height = msg.Height - 2
		v.ready = true
		if v.content != "" {
			v.viewport.SetContent(v.content)
		}
	case VSCodeFilesMsg:
		v.files = msg.Files
		if len(v.files) > 0 {
			path := v.files[0]
			return v, func() tea.Msg { return VSCodeFileSelectMsg{Path: path} }
		}
	case VSCodeContentMsg:
		v.content = msg.Content
		v.contentPath = msg.Path
		if v.ready {
			v.viewport.SetContent(msg.Content)
			v.viewport.GotoTop()
		}
	case tea.KeyMsg:
		if v.focusContent {
			switch msg.String() {
			case "h", "left":
				v.focusContent = false
				return v, nil
			}
			var cmd tea.Cmd
			v.viewport, cmd = v.viewport.Update(msg)
			return v, cmd
		}

		switch msg.String() {
		case "j", "down":
			if v.selected < len(v.files)-1 {
				v.selected++
				v.ensureListVisible()
				path := v.files[v.selected]
				return v, func() tea.Msg { return VSCodeFileSelectMsg{Path: path} }
			}
		case "k", "up":
			if v.selected > 0 {
				v.selected--
				v.ensureListVisible()
				path := v.files[v.selected]
				return v, func() tea.Msg { return VSCodeFileSelectMsg{Path: path} }
			}
		case "l", "right", "enter":
			v.focusContent = true
			return v, nil
		}
	}
	return v, nil
}

func (v *VSCode) ensureListVisible() {
	visible := v.height - 3
	if visible <= 0 {
		visible = 10
	}
	if v.selected < v.listOffset {
		v.listOffset = v.selected
	}
	if v.selected >= v.listOffset+visible {
		v.listOffset = v.selected - visible + 1
	}
}

func (v *VSCode) View() string {
	if len(v.files) == 0 {
		return vsMuted.Render("  Loading...")
	}

	listWidth := 32
	if v.width > 0 && v.width < 80 {
		listWidth = v.width / 3
	}

	// File list pane.
	var list strings.Builder
	list.WriteString(vsHeader.Render("Files") + "\n\n")
	visible := v.height - 3
	if visible <= 0 {
		visible = len(v.files)
	}
	end := min(v.listOffset+visible, len(v.files))
	for i := v.listOffset; i < end; i++ {
		cursor := "  "
		if i == v.selected {
			cursor = vsCursor.Render("> ")
		}
		name := filepath.Base(v.files[i])
		dir := filepath.Dir(v.files[i])
		style := vsNormal
		if i == v.selected {
			style = vsActive
		}
		entry := style.Render(name)
		if dir != "." {
			entry += vsMuted.Render(" " + dir)
		}
		list.WriteString(cursor + entry + "\n")
	}
	listPane := lipgloss.NewStyle().Width(listWidth).Render(list.String())

	// Content pane.
	var content strings.Builder
	if v.contentPath != "" {
		content.WriteString(vsMuted.Render(v.contentPath) + "\n\n")
	}
	if v.ready {
		content.WriteString(v.viewport.View())
	}

	// Border indicator for focus.
	contentStyle := lipgloss.NewStyle()
	if v.focusContent {
		contentStyle = contentStyle.BorderLeft(true).BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#9ccfd8"))
	} else {
		contentStyle = contentStyle.PaddingLeft(1)
	}

	contentPane := contentStyle.Render(content.String())

	return lipgloss.JoinHorizontal(lipgloss.Top, listPane, contentPane)
}
