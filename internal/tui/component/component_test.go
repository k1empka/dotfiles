package component

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestHelp_Toggle(t *testing.T) {
	h := NewHelp()
	if h.Visible() {
		t.Error("expected hidden initially")
	}
	h.Toggle()
	if !h.Visible() {
		t.Error("expected visible after toggle")
	}
	h.Toggle()
	if h.Visible() {
		t.Error("expected hidden after second toggle")
	}
}

func TestHelp_View(t *testing.T) {
	h := NewHelp()
	h.SetSize(80, 40)
	if h.View() != "" {
		t.Error("expected empty when hidden")
	}
	h.Toggle()
	view := h.View()
	if !strings.Contains(view, "Keyboard Shortcuts") {
		t.Error("expected title in view")
	}
}

func TestHelp_DismissEsc(t *testing.T) {
	h := NewHelp()
	h.Toggle()
	h, _ = h.Update(tea.KeyMsg{Type: tea.KeyEscape})
	if h.Visible() {
		t.Error("expected hidden after esc")
	}
}

func TestHelp_DismissQ(t *testing.T) {
	h := NewHelp()
	h.Toggle()
	h, _ = h.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if h.Visible() {
		t.Error("expected hidden after q")
	}
}
