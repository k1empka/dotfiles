package panel

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestOverview_Title(t *testing.T) {
	o := NewOverview()
	if o.Title() != "Overview" {
		t.Errorf("expected Overview, got %s", o.Title())
	}
}

func TestOverview_ConfigData(t *testing.T) {
	o := NewOverview()
	msg := ConfigDataMsg{Name: "Alice", Email: "alice@test.com", Theme: "rose-pine", OS: "darwin", Version: "2.40.0"}
	model, _ := o.Update(msg)
	o = model.(*Overview)
	view := o.View()
	if !strContains(view, "Alice") {
		t.Error("expected name in view")
	}
	if !strContains(view, "alice@test.com") {
		t.Error("expected email in view")
	}
}

func TestOverview_EditMode(t *testing.T) {
	o := NewOverview()
	o.name = "Alice"
	o.email = "alice@test.com"
	o.theme = "rose-pine"

	// Enter edit mode.
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("e")}
	model, _ := o.Update(msg)
	o = model.(*Overview)
	if !o.editing {
		t.Error("expected editing mode")
	}

	// Cancel with esc.
	msg = tea.KeyMsg{Type: tea.KeyEscape}
	model, _ = o.Update(msg)
	o = model.(*Overview)
	if o.editing {
		t.Error("expected editing cancelled")
	}
}

func TestShell_SubTabs(t *testing.T) {
	s := NewShell()
	s.width = 120
	s.height = 40

	// Set up size.
	model, _ := s.Update(SetSizeMsg{Width: 120, Height: 40})
	s = model.(*Shell)

	// Load content.
	model, _ = s.Update(ShellContentMsg{Index: 0, Content: "# zsh config"})
	s = model.(*Shell)
	model, _ = s.Update(ShellContentMsg{Index: 1, Content: "# bash config"})
	s = model.(*Shell)

	if s.subTab != 0 {
		t.Errorf("expected sub-tab 0, got %d", s.subTab)
	}

	// Switch to bash.
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")}
	model, _ = s.Update(msg)
	s = model.(*Shell)
	if s.subTab != 1 {
		t.Errorf("expected sub-tab 1, got %d", s.subTab)
	}

	// Switch back to zsh.
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")}
	model, _ = s.Update(msg)
	s = model.(*Shell)
	if s.subTab != 0 {
		t.Errorf("expected sub-tab 0, got %d", s.subTab)
	}
}

func TestGit_ContentMsg(t *testing.T) {
	g := NewGit()
	model, _ := g.Update(SetSizeMsg{Width: 80, Height: 30})
	g = model.(*Git)
	model, _ = g.Update(GitContentMsg{Content: "[user]\n  name = Alice"})
	g = model.(*Git)
	if g.content != "[user]\n  name = Alice" {
		t.Errorf("unexpected content: %q", g.content)
	}
}

func TestThemes_Navigation(t *testing.T) {
	th := NewThemes()
	th.themes = []string{"rose-pine", "catppuccin", "tokyo-night"}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
	model, _ := th.Update(msg)
	th = model.(*Themes)
	if th.selected != 1 {
		t.Errorf("expected selected 1, got %d", th.selected)
	}

	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")}
	model, _ = th.Update(msg)
	th = model.(*Themes)
	if th.selected != 0 {
		t.Errorf("expected selected 0, got %d", th.selected)
	}
}

func TestThemes_Confirm(t *testing.T) {
	th := NewThemes()
	th.themes = []string{"rose-pine", "catppuccin", "tokyo-night"}

	// Press enter to start confirm.
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	model, _ := th.Update(msg)
	th = model.(*Themes)
	if !th.confirming {
		t.Error("expected confirming mode")
	}

	// Press n to cancel.
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")}
	model, _ = th.Update(msg)
	th = model.(*Themes)
	if th.confirming {
		t.Error("expected confirming cancelled")
	}
}

func TestNeovim_Title(t *testing.T) {
	n := NewNeovim()
	if n.Title() != "Neovim" {
		t.Errorf("expected Neovim, got %s", n.Title())
	}
}

func TestNeovim_FileList(t *testing.T) {
	n := NewNeovim()
	model, _ := n.Update(NvimFilesMsg{Files: []string{"init.lua", "lua/plugins/user.lua"}})
	n = model.(*Neovim)
	if len(n.files) != 2 {
		t.Errorf("expected 2 files, got %d", len(n.files))
	}
}

func TestChezmoi_Title(t *testing.T) {
	c := NewChezmoi()
	if c.Title() != "Chezmoi" {
		t.Errorf("expected Chezmoi, got %s", c.Title())
	}
}

func TestChezmoi_StatusCommand(t *testing.T) {
	c := NewChezmoi()
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")}
	model, cmd := c.Update(msg)
	c = model.(*Chezmoi)
	if !c.loading {
		t.Error("expected loading state")
	}
	if cmd == nil {
		t.Error("expected command to be returned")
	}
}

func TestChezmoi_ApplyConfirm(t *testing.T) {
	c := NewChezmoi()

	// Press a to enter confirm.
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")}
	model, _ := c.Update(msg)
	c = model.(*Chezmoi)
	if !c.confirming {
		t.Error("expected confirming state")
	}

	// Press n to cancel.
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")}
	model, _ = c.Update(msg)
	c = model.(*Chezmoi)
	if c.confirming {
		t.Error("expected confirm cancelled")
	}
}

func TestChezmoi_Output(t *testing.T) {
	c := NewChezmoi()
	model, _ := c.Update(SetSizeMsg{Width: 80, Height: 30})
	c = model.(*Chezmoi)

	model, _ = c.Update(ChezmoiOutputMsg{Mode: "status", Output: "RA file.txt"})
	c = model.(*Chezmoi)
	if c.output != "RA file.txt" {
		t.Errorf("unexpected output: %q", c.output)
	}
	if c.loading {
		t.Error("expected loading to be false after output")
	}
}

func TestOverview_ShortHelp(t *testing.T) {
	o := NewOverview()
	help := o.ShortHelp()
	if !strContains(help, "edit") {
		t.Error("expected edit in help")
	}
}

func TestOverview_ViewEmpty(t *testing.T) {
	o := NewOverview()
	view := o.View()
	if !strContains(view, "not set") {
		t.Error("expected (not set) placeholder")
	}
}

func TestOverview_EditSubmit(t *testing.T) {
	o := NewOverview()
	o.name = "Alice"
	o.email = "alice@test.com"
	o.theme = "rose-pine"

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("e")}
	model, _ := o.Update(msg)
	o = model.(*Overview)
	if !o.editing {
		t.Fatal("expected editing mode")
	}

	// Submit.
	msg = tea.KeyMsg{Type: tea.KeyEnter}
	model, cmd := o.Update(msg)
	o = model.(*Overview)
	if o.editing {
		t.Error("expected editing to end after submit")
	}
	if cmd == nil {
		t.Error("expected command on submit")
	}
}

func TestOverview_ConfigUpdated(t *testing.T) {
	o := NewOverview()
	model, _ := o.Update(ConfigUpdatedMsg{Err: nil})
	o = model.(*Overview)
	if !strContains(o.View(), "updated") {
		t.Error("expected success message")
	}
}

func TestShell_Title(t *testing.T) {
	s := NewShell()
	if s.Title() != "Shell" {
		t.Errorf("expected Shell, got %s", s.Title())
	}
}

func TestShell_View(t *testing.T) {
	s := NewShell()
	model, _ := s.Update(SetSizeMsg{Width: 80, Height: 30})
	s = model.(*Shell)
	model, _ = s.Update(ShellContentMsg{Index: 0, Content: "# zsh"})
	s = model.(*Shell)
	view := s.View()
	if view == "" {
		t.Error("expected non-empty view")
	}
}

func TestGit_Title(t *testing.T) {
	g := NewGit()
	if g.Title() != "Git" {
		t.Errorf("expected Git, got %s", g.Title())
	}
}

func TestGit_ViewEmpty(t *testing.T) {
	g := NewGit()
	view := g.View()
	if !strContains(view, "Loading") {
		t.Error("expected Loading message")
	}
}

func TestThemes_Title(t *testing.T) {
	th := NewThemes()
	if th.Title() != "Themes" {
		t.Errorf("expected Themes, got %s", th.Title())
	}
}

func TestThemes_View(t *testing.T) {
	th := NewThemes()
	th.themes = []string{"rose-pine", "catppuccin"}
	th.active = "rose-pine"
	view := th.View()
	if !strContains(view, "rose-pine") {
		t.Error("expected rose-pine in view")
	}
	if !strContains(view, "active") {
		t.Error("expected (active) marker")
	}
}

func TestThemes_Colors(t *testing.T) {
	th := NewThemes()
	th.themes = []string{"rose-pine"}
	model, _ := th.Update(ThemeColorsMsg{Theme: "rose-pine", Colors: []string{"#e0def4", "#191724"}})
	th = model.(*Themes)
	if len(th.colors["rose-pine"]) != 2 {
		t.Errorf("expected 2 colors, got %d", len(th.colors["rose-pine"]))
	}
}

func TestThemes_ConfirmYes(t *testing.T) {
	th := NewThemes()
	th.themes = []string{"rose-pine", "catppuccin"}

	// Enter confirm.
	model, _ := th.Update(tea.KeyMsg{Type: tea.KeyEnter})
	th = model.(*Themes)
	if !th.confirming {
		t.Fatal("expected confirming")
	}

	// Confirm.
	model, cmd := th.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("y")})
	th = model.(*Themes)
	if th.confirming {
		t.Error("expected confirming cleared")
	}
	if cmd == nil {
		t.Error("expected command")
	}
}

func TestNeovim_ShortHelp(t *testing.T) {
	n := NewNeovim()
	help := n.ShortHelp()
	if !strContains(help, "navigate") {
		t.Error("expected navigate in help")
	}
	n.focusContent = true
	help = n.ShortHelp()
	if !strContains(help, "file list") {
		t.Error("expected file list in content-focus help")
	}
}

func TestNeovim_View(t *testing.T) {
	n := NewNeovim()
	model, _ := n.Update(SetSizeMsg{Width: 100, Height: 30})
	n = model.(*Neovim)
	model, _ = n.Update(NvimFilesMsg{Files: []string{"init.lua", "lua/lazy_setup.lua"}})
	n = model.(*Neovim)
	model, _ = n.Update(NvimContentMsg{Path: "init.lua", Content: "-- neovim init"})
	n = model.(*Neovim)
	view := n.View()
	if !strContains(view, "init.lua") {
		t.Error("expected init.lua in view")
	}
}

func TestNeovim_FocusSwitch(t *testing.T) {
	n := NewNeovim()
	n.files = []string{"a.lua", "b.lua"}
	model, _ := n.Update(SetSizeMsg{Width: 100, Height: 30})
	n = model.(*Neovim)

	// Focus content.
	model, _ = n.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	n = model.(*Neovim)
	if !n.focusContent {
		t.Error("expected content focused")
	}

	// Focus back to list.
	model, _ = n.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")})
	n = model.(*Neovim)
	if n.focusContent {
		t.Error("expected list focused")
	}
}

func TestChezmoi_ShortHelp(t *testing.T) {
	c := NewChezmoi()
	help := c.ShortHelp()
	if !strContains(help, "status") {
		t.Error("expected status in help")
	}
}

func TestChezmoi_DiffMode(t *testing.T) {
	c := NewChezmoi()
	model, _ := c.Update(SetSizeMsg{Width: 80, Height: 30})
	c = model.(*Chezmoi)

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("d")}
	model, _ = c.Update(msg)
	c = model.(*Chezmoi)
	if c.mode != "diff" {
		t.Errorf("expected diff mode, got %s", c.mode)
	}
}

func TestChezmoi_View(t *testing.T) {
	c := NewChezmoi()
	model, _ := c.Update(SetSizeMsg{Width: 80, Height: 30})
	c = model.(*Chezmoi)
	view := c.View()
	if !strContains(view, "Chezmoi Operations") {
		t.Error("expected header in view")
	}
}

func TestChezmoi_ConfirmView(t *testing.T) {
	c := NewChezmoi()
	model, _ := c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
	c = model.(*Chezmoi)
	view := c.View()
	if !strContains(view, "Apply") {
		t.Error("expected Apply question in view")
	}
}

func TestAlacritty_Title(t *testing.T) {
	a := NewAlacritty()
	if a.Title() != "Alacritty" {
		t.Errorf("expected Alacritty, got %s", a.Title())
	}
}

func TestAlacritty_ContentMsg(t *testing.T) {
	a := NewAlacritty()
	model, _ := a.Update(SetSizeMsg{Width: 80, Height: 30})
	a = model.(*Alacritty)
	model, _ = a.Update(AlacrittyContentMsg{Content: `import = ["~/.config/alacritty/themes/rose-pine.toml"]`})
	a = model.(*Alacritty)
	if a.content != `import = ["~/.config/alacritty/themes/rose-pine.toml"]` {
		t.Errorf("unexpected content: %q", a.content)
	}
}

func TestAlacritty_ViewEmpty(t *testing.T) {
	a := NewAlacritty()
	view := a.View()
	if !strContains(view, "Loading") {
		t.Error("expected Loading message")
	}
}

func TestAlacritty_ShortHelp(t *testing.T) {
	a := NewAlacritty()
	help := a.ShortHelp()
	if !strContains(help, "scroll") {
		t.Error("expected scroll in help")
	}
}

func TestAllPanels_ImplementInterface(t *testing.T) {
	panels := []Panel{
		NewOverview(),
		NewShell(),
		NewGit(),
		NewThemes(),
		NewNeovim(),
		NewAlacritty(),
		NewChezmoi(),
	}
	for _, p := range panels {
		if p.Title() == "" {
			t.Error("expected non-empty title")
		}
		if p.ShortHelp() == "" {
			t.Error("expected non-empty help")
		}
		_ = p.Init()
	}
}

func strContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
