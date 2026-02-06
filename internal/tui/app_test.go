package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/k1empka/dotfiles/internal/installer"
	"github.com/k1empka/dotfiles/internal/tui/panel"
)

func TestNewApp_HasEightPanels(t *testing.T) {
	app := NewApp(nil, nil)
	if len(app.panels) != 8 {
		t.Errorf("expected 8 panels, got %d", len(app.panels))
	}
}

func TestApp_TabSwitching(t *testing.T) {
	app := NewApp(nil, nil)
	app.width = 120
	app.height = 40

	tests := []struct {
		name    string
		key     string
		wantTab int
	}{
		{"tab 2", "2", 1},
		{"tab 3", "3", 2},
		{"tab 1", "1", 0},
		{"tab 4", "4", 3},
		{"tab 5", "5", 4},
		{"tab 6", "6", 5},
		{"tab 7", "7", 6},
		{"tab 8", "8", 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			model, _ := app.Update(msg)
			app = model.(App)
			if app.tabs.active != tt.wantTab {
				t.Errorf("expected tab %d, got %d", tt.wantTab, app.tabs.active)
			}
		})
	}
}

func TestApp_TabNext(t *testing.T) {
	app := NewApp(nil, nil)
	app.width = 120
	app.height = 40

	msg := tea.KeyMsg{Type: tea.KeyTab}
	model, _ := app.Update(msg)
	app = model.(App)
	if app.tabs.active != 1 {
		t.Errorf("expected tab 1 after Tab, got %d", app.tabs.active)
	}
}

func TestApp_TabPrev(t *testing.T) {
	app := NewApp(nil, nil)
	app.width = 120
	app.height = 40

	msg := tea.KeyMsg{Type: tea.KeyShiftTab}
	model, _ := app.Update(msg)
	app = model.(App)
	if app.tabs.active != 7 {
		t.Errorf("expected tab 7 after Shift+Tab from 0, got %d", app.tabs.active)
	}
}

func TestApp_Quit(t *testing.T) {
	app := NewApp(nil, nil)
	app.width = 120
	app.height = 40

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}
	_, cmd := app.Update(msg)
	if cmd == nil {
		t.Fatal("expected quit command")
	}
}

func TestApp_WindowResize(t *testing.T) {
	app := NewApp(nil, nil)
	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	model, _ := app.Update(msg)
	app = model.(App)
	if app.width != 100 || app.height != 50 {
		t.Errorf("expected 100x50, got %dx%d", app.width, app.height)
	}
}

func TestApp_ConfigLoaded(t *testing.T) {
	app := NewApp(nil, nil)
	app.width = 120
	app.height = 40

	// Broadcast size first so panels are ready.
	model, _ := app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	app = model.(App)

	msg := ConfigLoadedMsg{
		Name:    "Test User",
		Email:   "test@example.com",
		Theme:   "rose-pine",
		OS:      "darwin",
		Version: "2.40.0",
	}
	model, _ = app.Update(msg)
	app = model.(App)

	// Verify overview panel has the data by checking its view.
	view := app.panels[0].View()
	if !containsStr(view, "Test User") {
		t.Error("expected overview to show user name")
	}
}

func TestApp_HelpToggle(t *testing.T) {
	app := NewApp(nil, nil)
	app.width = 120
	app.height = 40

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")}
	model, _ := app.Update(msg)
	app = model.(App)
	if !app.help.Visible() {
		t.Error("expected help to be visible after ?")
	}

	// Press ? again to dismiss.
	model, _ = app.Update(msg)
	app = model.(App)
	if app.help.Visible() {
		t.Error("expected help to be hidden after second ?")
	}
}

func TestApp_ViewWithoutSize(t *testing.T) {
	app := NewApp(nil, nil)
	view := app.View()
	if view != "Loading..." {
		t.Errorf("expected Loading..., got %q", view)
	}
}

func TestApp_PanelTitles(t *testing.T) {
	app := NewApp(nil, nil)
	expected := []string{"Overview", "Shell", "Git", "Themes", "Neovim", "Alacritty", "Chezmoi", "Install"}
	for i, p := range app.panels {
		if p.Title() != expected[i] {
			t.Errorf("panel %d: expected title %q, got %q", i, expected[i], p.Title())
		}
	}
}

func TestApp_ErrorHandling(t *testing.T) {
	app := NewApp(nil, nil)
	app.width = 120
	app.height = 40

	msg := ConfigLoadedMsg{Err: errTest}
	model, _ := app.Update(msg)
	app = model.(App)
	if app.err == nil {
		t.Error("expected error to be set")
	}
}

func TestApp_EditSubmit(t *testing.T) {
	app := NewApp(nil, nil)
	app.width = 120
	app.height = 40

	msg := panel.EditSubmitMsg{Values: map[string]string{"name": "New Name"}}
	_, cmd := app.Update(msg)
	if cmd == nil {
		t.Error("expected command for config write")
	}
}

var errTest = &testError{}

type testError struct{}

func (e *testError) Error() string { return "test error" }

func TestApp_ViewWithSize(t *testing.T) {
	app := NewApp(nil, nil)
	model, _ := app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	app = model.(App)
	view := app.View()
	if view == "" {
		t.Error("expected non-empty view")
	}
	if !containsStr(view, "Overview") {
		t.Error("expected Overview tab in view")
	}
}

func TestApp_ViewWithHelp(t *testing.T) {
	app := NewApp(nil, nil)
	app.width = 120
	app.height = 40
	app.help.SetSize(120, 40)
	app.help.Toggle()

	view := app.View()
	if !containsStr(view, "Keyboard Shortcuts") {
		t.Error("expected help overlay")
	}
}

func TestApp_ViewEachTab(t *testing.T) {
	app := NewApp(nil, nil)
	model, _ := app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	app = model.(App)

	for i := range 8 {
		app.tabs.setActive(i)
		view := app.View()
		if view == "" {
			t.Errorf("tab %d: expected non-empty view", i)
		}
	}
}

func TestApp_RunChezmoiMsg(t *testing.T) {
	app := NewApp(nil, nil)
	app.width = 120
	app.height = 40

	msg := panel.RunChezmoiMsg{Mode: "status"}
	_, cmd := app.Update(msg)
	if cmd == nil {
		t.Error("expected command for chezmoi status")
	}
}

func TestApp_ThemeApply(t *testing.T) {
	app := NewApp(nil, nil)
	app.width = 120
	app.height = 40

	msg := panel.ThemeApplyMsg{Theme: "catppuccin"}
	_, cmd := app.Update(msg)
	if cmd == nil {
		t.Error("expected command for theme apply")
	}
}

func TestApp_BroadcastMessages(t *testing.T) {
	app := NewApp(nil, nil)
	model, _ := app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	app = model.(App)

	// Shell content should be broadcast.
	msg := panel.ShellContentMsg{Index: 0, Content: "# zsh"}
	model, _ = app.Update(msg)
	app = model.(App)

	// Git content.
	gitMsg := panel.GitContentMsg{Content: "[user]"}
	model, _ = app.Update(gitMsg)
	app = model.(App)

	// Alacritty content.
	alMsg := panel.AlacrittyContentMsg{Content: `import = ["~/.config/alacritty/themes/rose-pine.toml"]`}
	model, _ = app.Update(alMsg)
	app = model.(App)

	// Nvim files.
	nvimMsg := panel.NvimFilesMsg{Files: []string{"init.lua"}}
	_, _ = app.Update(nvimMsg)
}

func TestApp_NvimFileSelect(t *testing.T) {
	app := NewApp(nil, nil)
	msg := panel.NvimFileSelectMsg{Path: "init.lua"}
	_, cmd := app.Update(msg)
	if cmd == nil {
		t.Error("expected command for file load")
	}
}

func TestApp_CheckInstallMsg(t *testing.T) {
	app := NewApp(nil, nil)
	app.width = 120
	app.height = 40

	msg := panel.CheckInstallMsg{}
	_, cmd := app.Update(msg)
	if cmd == nil {
		t.Error("expected command for install check")
	}
}

func TestApp_InstallRequestMsg(t *testing.T) {
	app := NewApp(nil, nil)
	app.width = 120
	app.height = 40

	msg := panel.InstallRequestMsg{App: installer.App{Name: "Neovim", BinName: "nvim"}}
	_, cmd := app.Update(msg)
	if cmd == nil {
		t.Error("expected command for app install")
	}
}

func TestApp_InstallStatusBroadcast(t *testing.T) {
	app := NewApp(nil, nil)
	model, _ := app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	app = model.(App)

	msg := panel.InstallStatusMsg{Statuses: []installer.AppStatus{
		{App: installer.App{Name: "Neovim", BinName: "nvim"}, Status: installer.StatusInstalled},
	}}
	model, _ = app.Update(msg)
	app = model.(App)
	// Verify no panic during broadcast.
	if app.width != 120 {
		t.Error("expected width to remain 120")
	}
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
