package panel

import (
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	nvHeader = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#c4a7e7"))
	nvCursor = lipgloss.NewStyle().Foreground(lipgloss.Color("#9ccfd8"))
	nvActive = lipgloss.NewStyle().Foreground(lipgloss.Color("#e0def4")).Bold(true)
	nvNormal = lipgloss.NewStyle().Foreground(lipgloss.Color("#908caa"))
	nvMuted  = lipgloss.NewStyle().Foreground(lipgloss.Color("#6e6a86"))
)

// NvimFileSelectMsg requests a file content load from the app.
type NvimFileSelectMsg struct {
	Path string
}

// Neovim displays neovim config files in a split-pane browser.
type Neovim struct {
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

func NewNeovim() *Neovim { return &Neovim{} }

func (n *Neovim) Title() string { return "Neovim" }

func (n *Neovim) ShortHelp() string {
	if n.focusContent {
		return "h: file list  j/k: scroll  q: quit  ?: help"
	}
	return "j/k: navigate  l/enter: view file  q: quit  ?: help"
}

func (n *Neovim) Init() tea.Cmd { return nil }

func (n *Neovim) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SetSizeMsg:
		n.width = msg.Width
		n.height = msg.Height
		n.viewport.Width = msg.Width - 34
		n.viewport.Height = msg.Height - 2
		n.ready = true
		if n.content != "" {
			n.viewport.SetContent(n.content)
		}
	case NvimFilesMsg:
		n.files = msg.Files
		if len(n.files) > 0 {
			path := n.files[0]
			return n, func() tea.Msg { return NvimFileSelectMsg{Path: path} }
		}
	case NvimContentMsg:
		n.content = msg.Content
		n.contentPath = msg.Path
		if n.ready {
			n.viewport.SetContent(msg.Content)
			n.viewport.GotoTop()
		}
	case tea.KeyMsg:
		if n.focusContent {
			switch msg.String() {
			case "h", "left":
				n.focusContent = false
				return n, nil
			}
			var cmd tea.Cmd
			n.viewport, cmd = n.viewport.Update(msg)
			return n, cmd
		}

		switch msg.String() {
		case "j", "down":
			if n.selected < len(n.files)-1 {
				n.selected++
				n.ensureListVisible()
				path := n.files[n.selected]
				return n, func() tea.Msg { return NvimFileSelectMsg{Path: path} }
			}
		case "k", "up":
			if n.selected > 0 {
				n.selected--
				n.ensureListVisible()
				path := n.files[n.selected]
				return n, func() tea.Msg { return NvimFileSelectMsg{Path: path} }
			}
		case "l", "right", "enter":
			n.focusContent = true
			return n, nil
		}
	}
	return n, nil
}

func (n *Neovim) ensureListVisible() {
	visible := n.height - 3
	if visible <= 0 {
		visible = 10
	}
	if n.selected < n.listOffset {
		n.listOffset = n.selected
	}
	if n.selected >= n.listOffset+visible {
		n.listOffset = n.selected - visible + 1
	}
}

func (n *Neovim) View() string {
	if len(n.files) == 0 {
		return nvMuted.Render("  Loading...")
	}

	listWidth := 32
	if n.width > 0 && n.width < 80 {
		listWidth = n.width / 3
	}

	// File list pane.
	var list strings.Builder
	list.WriteString(nvHeader.Render("Files") + "\n\n")
	visible := n.height - 3
	if visible <= 0 {
		visible = len(n.files)
	}
	end := min(n.listOffset+visible, len(n.files))
	for i := n.listOffset; i < end; i++ {
		cursor := "  "
		if i == n.selected {
			cursor = nvCursor.Render("> ")
		}
		name := filepath.Base(n.files[i])
		dir := filepath.Dir(n.files[i])
		style := nvNormal
		if i == n.selected {
			style = nvActive
		}
		entry := style.Render(name)
		if dir != "." {
			entry += nvMuted.Render(" " + dir)
		}
		list.WriteString(cursor + entry + "\n")
	}
	listPane := lipgloss.NewStyle().Width(listWidth).Render(list.String())

	// Content pane.
	var content strings.Builder
	if n.contentPath != "" {
		content.WriteString(nvMuted.Render(n.contentPath) + "\n\n")
	}
	if n.ready {
		content.WriteString(n.viewport.View())
	}

	// Border indicator for focus.
	contentStyle := lipgloss.NewStyle()
	if n.focusContent {
		contentStyle = contentStyle.BorderLeft(true).BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#9ccfd8"))
	} else {
		contentStyle = contentStyle.PaddingLeft(1)
	}

	contentPane := contentStyle.Render(content.String())

	return lipgloss.JoinHorizontal(lipgloss.Top, listPane, contentPane)
}
