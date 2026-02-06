package component

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	selectorCursor   = lipgloss.NewStyle().Foreground(lipgloss.Color("#9ccfd8")).Render("> ")
	selectorNoCursor = "  "
	selectorActive   = lipgloss.NewStyle().Foreground(lipgloss.Color("#e0def4")).Bold(true)
	selectorNormal   = lipgloss.NewStyle().Foreground(lipgloss.Color("#908caa"))
)

// SelectorItem represents an item in the selector list.
type SelectorItem struct {
	Label string
	Value string
}

// ItemSelectedMsg is emitted when an item is selected.
type ItemSelectedMsg struct {
	Index int
	Item  SelectorItem
}

// Selector is a simple navigable list of items.
type Selector struct {
	items    []SelectorItem
	cursor   int
	width    int
	height   int
	offset   int
}

// NewSelector creates a new Selector.
func NewSelector(items []SelectorItem) Selector {
	return Selector{items: items}
}

// SetItems replaces the current item list.
func (s *Selector) SetItems(items []SelectorItem) {
	s.items = items
	if s.cursor >= len(items) {
		s.cursor = max(0, len(items)-1)
	}
}

// SetSize updates the selector dimensions.
func (s *Selector) SetSize(w, h int) {
	s.width = w
	s.height = h
}

// Selected returns the currently highlighted item index.
func (s Selector) Selected() int { return s.cursor }

// SelectedItem returns the currently highlighted item.
func (s Selector) SelectedItem() (SelectorItem, bool) {
	if s.cursor < 0 || s.cursor >= len(s.items) {
		return SelectorItem{}, false
	}
	return s.items[s.cursor], true
}

func (s Selector) Init() tea.Cmd { return nil }

func (s Selector) Update(msg tea.Msg) (Selector, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if s.cursor < len(s.items)-1 {
				s.cursor++
				s.ensureVisible()
			}
		case "k", "up":
			if s.cursor > 0 {
				s.cursor--
				s.ensureVisible()
			}
		case "g":
			s.cursor = 0
			s.offset = 0
		case "G":
			s.cursor = max(0, len(s.items)-1)
			s.ensureVisible()
		case "enter":
			if item, ok := s.SelectedItem(); ok {
				return s, func() tea.Msg {
					return ItemSelectedMsg{Index: s.cursor, Item: item}
				}
			}
		}
	}
	return s, nil
}

func (s *Selector) ensureVisible() {
	visible := s.height
	if visible <= 0 {
		visible = 10
	}
	if s.cursor < s.offset {
		s.offset = s.cursor
	}
	if s.cursor >= s.offset+visible {
		s.offset = s.cursor - visible + 1
	}
}

func (s Selector) View() string {
	if len(s.items) == 0 {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#6e6a86")).Render("  (empty)")
	}

	visible := s.height
	if visible <= 0 {
		visible = len(s.items)
	}

	var view string
	end := min(s.offset+visible, len(s.items))
	for i := s.offset; i < end; i++ {
		cursor := selectorNoCursor
		style := selectorNormal
		if i == s.cursor {
			cursor = selectorCursor
			style = selectorActive
		}
		view += cursor + style.Render(s.items[i].Label) + "\n"
	}
	return view
}
