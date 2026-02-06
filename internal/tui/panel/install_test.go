package panel

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/k1empka/dotfiles/internal/installer"
)

func newTestInstall() *Install {
	p := NewInstall()
	p.width = 120
	p.height = 40
	return p
}

func testStatuses() []installer.AppStatus {
	return []installer.AppStatus{
		{App: installer.App{Name: "Neovim", BinName: "nvim"}, Status: installer.StatusInstalled, Version: "NVIM v0.10.0"},
		{App: installer.App{Name: "Alacritty", BinName: "alacritty"}, Status: installer.StatusMissing},
		{App: installer.App{Name: "Git", BinName: "git"}, Status: installer.StatusInstalled, Version: "git version 2.43.0"},
	}
}

func TestInstall_Title(t *testing.T) {
	p := newTestInstall()
	if p.Title() != "Install" {
		t.Errorf("expected title 'Install', got %q", p.Title())
	}
}

func TestInstall_ShortHelp(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*Install)
		wantSubstr string
	}{
		{"normal", func(p *Install) {}, "j/k"},
		{"confirming", func(p *Install) { p.confirming = true }, "confirm"},
		{"installing", func(p *Install) { p.installing = 0 }, "installing"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newTestInstall()
			tt.setup(p)
			if !strings.Contains(p.ShortHelp(), tt.wantSubstr) {
				t.Errorf("expected help to contain %q, got %q", tt.wantSubstr, p.ShortHelp())
			}
		})
	}
}

func TestInstall_StatusMsg(t *testing.T) {
	p := newTestInstall()
	model, _ := p.Update(InstallStatusMsg{Statuses: testStatuses()})
	p = model.(*Install)

	if len(p.statuses) != 3 {
		t.Errorf("expected 3 statuses, got %d", len(p.statuses))
	}
	if !p.loaded {
		t.Error("expected loaded to be true")
	}
}

func TestInstall_Navigation(t *testing.T) {
	p := newTestInstall()
	model, _ := p.Update(InstallStatusMsg{Statuses: testStatuses()})
	p = model.(*Install)

	// Move down.
	model, _ = p.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	p = model.(*Install)
	if p.selected != 1 {
		t.Errorf("expected selected 1, got %d", p.selected)
	}

	// Move down again.
	model, _ = p.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	p = model.(*Install)
	if p.selected != 2 {
		t.Errorf("expected selected 2, got %d", p.selected)
	}

	// Move down at bottom (stays).
	model, _ = p.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	p = model.(*Install)
	if p.selected != 2 {
		t.Errorf("expected selected 2, got %d", p.selected)
	}

	// Move up.
	model, _ = p.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	p = model.(*Install)
	if p.selected != 1 {
		t.Errorf("expected selected 1, got %d", p.selected)
	}
}

func TestInstall_EnterOnMissing(t *testing.T) {
	p := newTestInstall()
	model, _ := p.Update(InstallStatusMsg{Statuses: testStatuses()})
	p = model.(*Install)

	// Move to Alacritty (missing).
	model, _ = p.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	p = model.(*Install)

	// Press enter.
	model, _ = p.Update(tea.KeyMsg{Type: tea.KeyEnter})
	p = model.(*Install)

	if !p.confirming {
		t.Error("expected confirming to be true")
	}
}

func TestInstall_EnterOnInstalled(t *testing.T) {
	p := newTestInstall()
	model, _ := p.Update(InstallStatusMsg{Statuses: testStatuses()})
	p = model.(*Install)

	// Selected is Neovim (installed). Press enter.
	model, _ = p.Update(tea.KeyMsg{Type: tea.KeyEnter})
	p = model.(*Install)

	if p.confirming {
		t.Error("expected confirming to be false for installed app")
	}
}

func TestInstall_ConfirmYes(t *testing.T) {
	p := newTestInstall()
	model, _ := p.Update(InstallStatusMsg{Statuses: testStatuses()})
	p = model.(*Install)

	// Select missing app and confirm.
	model, _ = p.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	p = model.(*Install)
	model, _ = p.Update(tea.KeyMsg{Type: tea.KeyEnter})
	p = model.(*Install)
	model, cmd := p.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("y")})
	p = model.(*Install)

	if p.confirming {
		t.Error("expected confirming to be false after y")
	}
	if p.installing != 1 {
		t.Errorf("expected installing index 1, got %d", p.installing)
	}
	if cmd == nil {
		t.Error("expected command for install")
	}
}

func TestInstall_ConfirmNo(t *testing.T) {
	p := newTestInstall()
	model, _ := p.Update(InstallStatusMsg{Statuses: testStatuses()})
	p = model.(*Install)

	// Select missing app and cancel.
	model, _ = p.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	p = model.(*Install)
	model, _ = p.Update(tea.KeyMsg{Type: tea.KeyEnter})
	p = model.(*Install)
	model, _ = p.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	p = model.(*Install)

	if p.confirming {
		t.Error("expected confirming to be false after n")
	}
	if p.installing != -1 {
		t.Errorf("expected installing -1, got %d", p.installing)
	}
}

func TestInstall_Refresh(t *testing.T) {
	p := newTestInstall()
	model, _ := p.Update(InstallStatusMsg{Statuses: testStatuses()})
	p = model.(*Install)

	model, cmd := p.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")})
	p = model.(*Install)

	if !p.loading {
		t.Error("expected loading to be true after r")
	}
	if cmd == nil {
		t.Error("expected command for refresh")
	}
}

func TestInstall_ResultSuccess(t *testing.T) {
	p := newTestInstall()
	model, _ := p.Update(InstallStatusMsg{Statuses: testStatuses()})
	p = model.(*Install)

	model, _ = p.Update(InstallResultMsg{
		App: installer.App{Name: "Alacritty", BinName: "alacritty"},
		Err: nil,
	})
	p = model.(*Install)

	if p.statuses[1].Status != installer.StatusInstalled {
		t.Errorf("expected StatusInstalled, got %d", p.statuses[1].Status)
	}
}

func TestInstall_ResultError(t *testing.T) {
	p := newTestInstall()
	model, _ := p.Update(InstallStatusMsg{Statuses: testStatuses()})
	p = model.(*Install)

	testErr := &testError{}
	model, _ = p.Update(InstallResultMsg{
		App: installer.App{Name: "Alacritty", BinName: "alacritty"},
		Err: testErr,
	})
	p = model.(*Install)

	if p.statuses[1].Status != installer.StatusError {
		t.Errorf("expected StatusError, got %d", p.statuses[1].Status)
	}
}

func TestInstall_View_Loading(t *testing.T) {
	p := newTestInstall()
	view := p.View()
	if !strings.Contains(view, "Loading") {
		t.Error("expected Loading in view before data arrives")
	}
}

func TestInstall_View_Indicators(t *testing.T) {
	p := newTestInstall()
	model, _ := p.Update(InstallStatusMsg{Statuses: testStatuses()})
	p = model.(*Install)
	view := p.View()

	if !strings.Contains(view, "Neovim") {
		t.Error("expected Neovim in view")
	}
	if !strings.Contains(view, "Alacritty") {
		t.Error("expected Alacritty in view")
	}
	if !strings.Contains(view, "Application Status") {
		t.Error("expected header in view")
	}
}

type testError struct{}

func (e *testError) Error() string { return "test error" }
