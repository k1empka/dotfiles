package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// tabBar renders a horizontal tab bar.
type tabBar struct {
	titles []string
	active int
}

func newTabBar(titles []string) tabBar {
	return tabBar{titles: titles}
}

func (t *tabBar) setActive(i int) {
	if i >= 0 && i < len(t.titles) {
		t.active = i
	}
}

func (t *tabBar) next() {
	t.active = (t.active + 1) % len(t.titles)
}

func (t *tabBar) prev() {
	t.active = (t.active - 1 + len(t.titles)) % len(t.titles)
}

func (t tabBar) View(width int) string {
	var tabs []string
	for i, title := range t.titles {
		num := lipgloss.NewStyle().Foreground(colorMuted).Render(string(rune('1'+i)) + " ")
		if i == t.active {
			num = lipgloss.NewStyle().Foreground(colorFoam).Render(string(rune('1'+i)) + " ")
			tabs = append(tabs, activeTabStyle.Render(num+title))
		} else {
			tabs = append(tabs, inactiveTabStyle.Render(num+title))
		}
	}
	row := strings.Join(tabs, "")
	gap := width - lipgloss.Width(row)
	if gap > 0 {
		row += lipgloss.NewStyle().Background(colorBase).Render(strings.Repeat(" ", gap))
	}
	return row
}
