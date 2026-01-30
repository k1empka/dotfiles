package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var titleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#e0def4")).
	Background(lipgloss.Color("#191724")).
	Padding(0, 1)

// App is the root model for the dotfiles TUI.
type App struct {
	width  int
	height int
}

// NewApp creates a new App model.
func NewApp() App {
	return App{}
}

func (a App) Init() tea.Cmd {
	return nil
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return a, tea.Quit
		}
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
	}
	return a, nil
}

func (a App) View() string {
	title := titleStyle.Render("dotfiles-tui")
	help := "Press q to quit"
	return title + "\n\n  Dotfiles manager — coming soon.\n\n  " + help + "\n"
}
