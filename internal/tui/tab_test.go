package tui

import "testing"

func TestTabBar_Navigation(t *testing.T) {
	tb := newTabBar([]string{"A", "B", "C"})

	tb.next()
	if tb.active != 1 {
		t.Errorf("expected 1, got %d", tb.active)
	}

	tb.next()
	if tb.active != 2 {
		t.Errorf("expected 2, got %d", tb.active)
	}

	tb.next()
	if tb.active != 0 {
		t.Errorf("expected 0 (wrap), got %d", tb.active)
	}

	tb.prev()
	if tb.active != 2 {
		t.Errorf("expected 2 (wrap back), got %d", tb.active)
	}
}

func TestTabBar_SetActive(t *testing.T) {
	tb := newTabBar([]string{"A", "B", "C"})

	tb.setActive(2)
	if tb.active != 2 {
		t.Errorf("expected 2, got %d", tb.active)
	}

	// Out of range should be ignored.
	tb.setActive(-1)
	if tb.active != 2 {
		t.Errorf("expected 2 unchanged, got %d", tb.active)
	}

	tb.setActive(5)
	if tb.active != 2 {
		t.Errorf("expected 2 unchanged, got %d", tb.active)
	}
}

func TestTabBar_View(t *testing.T) {
	tb := newTabBar([]string{"Overview", "Shell", "Git"})
	view := tb.View(80)
	if view == "" {
		t.Error("expected non-empty tab bar view")
	}
}
