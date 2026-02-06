package component

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// --- Viewer tests ---

func TestViewer_SetContent(t *testing.T) {
	v := NewViewer(false)
	v.SetSize(80, 20)
	v.SetContent("Hello World")
	view := v.View()
	if view == "" {
		t.Error("expected non-empty view")
	}
}

func TestViewer_LineNumbers(t *testing.T) {
	v := NewViewer(true)
	v.SetSize(80, 20)
	v.SetContent("line one\nline two")
	view := v.View()
	if !strContains(view, "1") {
		t.Error("expected line number 1 in view")
	}
}

func TestViewer_NotReady(t *testing.T) {
	v := NewViewer(false)
	view := v.View()
	if view != "" {
		t.Error("expected empty view when not ready")
	}
}

func TestViewer_ScrollPercent(t *testing.T) {
	v := NewViewer(false)
	v.SetSize(80, 5)
	v.SetContent("1\n2\n3\n4\n5\n6\n7\n8\n9\n10")
	_ = v.ScrollPercent()
}

func TestViewer_Update(t *testing.T) {
	v := NewViewer(false)
	v.SetSize(80, 20)
	v.SetContent("test content")
	v, _ = v.Update(tea.KeyMsg{Type: tea.KeyDown})
	_ = v.View()
}

// --- Selector tests ---

func TestSelector_Navigation(t *testing.T) {
	items := []SelectorItem{
		{Label: "one", Value: "1"},
		{Label: "two", Value: "2"},
		{Label: "three", Value: "3"},
	}
	s := NewSelector(items)
	s.SetSize(40, 10)

	// Move down.
	s, _ = s.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	if s.Selected() != 1 {
		t.Errorf("expected 1, got %d", s.Selected())
	}

	// Move up.
	s, _ = s.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	if s.Selected() != 0 {
		t.Errorf("expected 0, got %d", s.Selected())
	}

	// Jump to bottom.
	s, _ = s.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("G")})
	if s.Selected() != 2 {
		t.Errorf("expected 2, got %d", s.Selected())
	}

	// Jump to top.
	s, _ = s.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("g")})
	if s.Selected() != 0 {
		t.Errorf("expected 0, got %d", s.Selected())
	}
}

func TestSelector_Enter(t *testing.T) {
	items := []SelectorItem{{Label: "item", Value: "val"}}
	s := NewSelector(items)
	s.SetSize(40, 10)

	_, cmd := s.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Error("expected command on enter")
	}
}

func TestSelector_Empty(t *testing.T) {
	s := NewSelector(nil)
	view := s.View()
	if !strContains(view, "empty") {
		t.Error("expected empty indicator")
	}
}

func TestSelector_SetItems(t *testing.T) {
	s := NewSelector(nil)
	s.SetItems([]SelectorItem{{Label: "new", Value: "1"}})
	if s.Selected() != 0 {
		t.Errorf("expected 0, got %d", s.Selected())
	}
}

func TestSelector_SelectedItem(t *testing.T) {
	s := NewSelector([]SelectorItem{{Label: "one", Value: "1"}})
	item, ok := s.SelectedItem()
	if !ok {
		t.Error("expected ok")
	}
	if item.Label != "one" {
		t.Errorf("expected one, got %s", item.Label)
	}
}

func TestSelector_View(t *testing.T) {
	items := []SelectorItem{{Label: "a"}, {Label: "b"}}
	s := NewSelector(items)
	s.SetSize(40, 10)
	view := s.View()
	if !strContains(view, "a") || !strContains(view, "b") {
		t.Error("expected items in view")
	}
}

// --- Confirm tests ---

func TestConfirm_ShowHide(t *testing.T) {
	c := NewConfirm()
	if c.Visible() {
		t.Error("expected hidden initially")
	}
	c.Show("Delete?", "delete")
	if !c.Visible() {
		t.Error("expected visible after Show")
	}
	c.Hide()
	if c.Visible() {
		t.Error("expected hidden after Hide")
	}
}

func TestConfirm_Yes(t *testing.T) {
	c := NewConfirm()
	c.Show("Proceed?", "test")

	c, cmd := c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("y")})
	if c.Visible() {
		t.Error("expected hidden after y")
	}
	if cmd == nil {
		t.Error("expected command")
	}
	msg := cmd()
	cm, ok := msg.(ConfirmMsg)
	if !ok {
		t.Fatal("expected ConfirmMsg")
	}
	if !cm.Confirmed || cm.Tag != "test" {
		t.Errorf("unexpected: confirmed=%v tag=%s", cm.Confirmed, cm.Tag)
	}
}

func TestConfirm_No(t *testing.T) {
	c := NewConfirm()
	c.Show("Proceed?", "test")

	c, cmd := c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	if c.Visible() {
		t.Error("expected hidden after n")
	}
	if cmd == nil {
		t.Error("expected command")
	}
	msg := cmd()
	cm := msg.(ConfirmMsg)
	if cm.Confirmed {
		t.Error("expected not confirmed")
	}
}

func TestConfirm_Hidden(t *testing.T) {
	c := NewConfirm()
	c, _ = c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("y")})
	if c.Visible() {
		t.Error("should stay hidden when not shown")
	}
}

func TestConfirm_View(t *testing.T) {
	c := NewConfirm()
	if c.View() != "" {
		t.Error("expected empty view when hidden")
	}
	c.Show("Proceed?", "test")
	view := c.View()
	if !strContains(view, "Proceed?") {
		t.Error("expected message in view")
	}
}

// --- DiffView tests ---

func TestDiffView_Colorize(t *testing.T) {
	d := NewDiffView()
	d.SetSize(80, 20)
	d.SetContent("--- a/file\n+++ b/file\n@@ -1,3 +1,3 @@\n-old\n+new\n same")
	view := d.View()
	if view == "" {
		t.Error("expected non-empty view")
	}
}

func TestDiffView_NotReady(t *testing.T) {
	d := NewDiffView()
	if d.View() != "" {
		t.Error("expected empty view when not ready")
	}
}

func TestDiffView_Update(t *testing.T) {
	d := NewDiffView()
	d.SetSize(80, 20)
	d.SetContent("test")
	d, _ = d.Update(tea.KeyMsg{Type: tea.KeyDown})
	_ = d.View()
}

// --- ErrorDisplay tests ---

func TestErrorDisplay_ShowDismiss(t *testing.T) {
	e := NewErrorDisplay()
	if e.Visible() {
		t.Error("expected hidden initially")
	}
	cmd := e.Show("something failed")
	if !e.Visible() {
		t.Error("expected visible after Show")
	}
	if cmd == nil {
		t.Error("expected tick command")
	}
	view := e.View()
	if !strContains(view, "something failed") {
		t.Error("expected message in view")
	}

	// Simulate clear.
	e, _ = e.Update(clearErrorMsg{})
	if e.Visible() {
		t.Error("expected hidden after clear")
	}
}

func TestErrorDisplay_HiddenView(t *testing.T) {
	e := NewErrorDisplay()
	if e.View() != "" {
		t.Error("expected empty view when hidden")
	}
}

// --- Help tests ---

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
	if !strContains(view, "Keyboard Shortcuts") {
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

// --- StatusBar tests ---

func TestStatusBar_Render(t *testing.T) {
	sb := NewStatusBar()
	sb.SetWidth(80)
	result := sb.Render("help text", nil)
	if !strContains(result, "help text") {
		t.Error("expected help text in render")
	}
}

func TestStatusBar_RenderWithError(t *testing.T) {
	sb := NewStatusBar()
	sb.SetWidth(80)
	result := sb.Render("help", &testErr{})
	// In non-TTY mode, lipgloss may strip styles; just check it returns something.
	if result == "" {
		t.Error("expected non-empty render")
	}
}

func TestFormatHelp(t *testing.T) {
	result := FormatHelp("q", "quit", "?", "help")
	if result == "" {
		t.Error("expected non-empty result")
	}
}

type testErr struct{}

func (e *testErr) Error() string { return "test error" }

func strContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
