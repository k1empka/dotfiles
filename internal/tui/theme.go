package tui

import "github.com/charmbracelet/lipgloss"

// Rose Pine palette.
var (
	colorBase    = lipgloss.Color("#191724")
	colorSurface = lipgloss.Color("#1f1d2e")
	colorOverlay = lipgloss.Color("#26233a")
	colorMuted   = lipgloss.Color("#6e6a86")
	colorSubtle  = lipgloss.Color("#908caa")
	colorText    = lipgloss.Color("#e0def4")
	colorLove = lipgloss.Color("#eb6f92")
	colorFoam = lipgloss.Color("#9ccfd8")
)

// Shared styles used across the TUI.
var (
	activeTabStyle = lipgloss.NewStyle().
			Background(colorOverlay).
			Foreground(colorText).
			Bold(true).
			Padding(0, 1)

	inactiveTabStyle = lipgloss.NewStyle().
				Background(colorBase).
				Foreground(colorMuted).
				Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
			Background(colorSurface).
			Foreground(colorSubtle).
			Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
			Foreground(colorLove)
)

