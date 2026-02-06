package installer

import (
	"fmt"
	"strings"
	"testing"
)

type mockRunner struct {
	lookPath map[string]error
	runs     map[string]runResult
}

type runResult struct {
	output []byte
	err    error
}

func (m *mockRunner) LookPath(file string) (string, error) {
	if err, ok := m.lookPath[file]; ok {
		return "/usr/bin/" + file, err
	}
	return "", fmt.Errorf("not found: %s", file)
}

func (m *mockRunner) Run(name string, args ...string) ([]byte, error) {
	key := name + " " + strings.Join(args, " ")
	key = strings.TrimSpace(key)
	if r, ok := m.runs[key]; ok {
		return r.output, r.err
	}
	return nil, fmt.Errorf("unexpected command: %s", key)
}

func TestCheck_Installed(t *testing.T) {
	runner := &mockRunner{
		lookPath: map[string]error{"nvim": nil},
		runs: map[string]runResult{
			"nvim --version": {output: []byte("NVIM v0.10.0\nBuild type: Release")},
		},
	}
	inst := NewInstaller(runner, "darwin")

	result := inst.Check(App{Name: "Neovim", BinName: "nvim"})

	if result.Status != StatusInstalled {
		t.Errorf("expected StatusInstalled, got %d", result.Status)
	}
	if result.Version != "NVIM v0.10.0" {
		t.Errorf("expected version 'NVIM v0.10.0', got %q", result.Version)
	}
}

func TestCheck_Missing(t *testing.T) {
	runner := &mockRunner{
		lookPath: map[string]error{},
		runs:     map[string]runResult{},
	}
	inst := NewInstaller(runner, "darwin")

	result := inst.Check(App{Name: "Neovim", BinName: "nvim"})

	if result.Status != StatusMissing {
		t.Errorf("expected StatusMissing, got %d", result.Status)
	}
}

func TestCheck_InstalledNoVersion(t *testing.T) {
	runner := &mockRunner{
		lookPath: map[string]error{"nvim": nil},
		runs:     map[string]runResult{},
	}
	inst := NewInstaller(runner, "darwin")

	result := inst.Check(App{Name: "Neovim", BinName: "nvim"})

	if result.Status != StatusInstalled {
		t.Errorf("expected StatusInstalled, got %d", result.Status)
	}
	if result.Version != "" {
		t.Errorf("expected empty version, got %q", result.Version)
	}
}

func TestCheckAll(t *testing.T) {
	runner := &mockRunner{
		lookPath: map[string]error{
			"nvim": nil,
			"git":  nil,
		},
		runs: map[string]runResult{
			"nvim --version": {output: []byte("NVIM v0.10.0")},
			"git --version":  {output: []byte("git version 2.43.0")},
		},
	}
	inst := NewInstaller(runner, "darwin")

	results := inst.CheckAll()

	if len(results) != len(Apps()) {
		t.Errorf("expected %d results, got %d", len(Apps()), len(results))
	}

	// Neovim and Git installed, others missing.
	installed := 0
	missing := 0
	for _, r := range results {
		switch r.Status {
		case StatusInstalled:
			installed++
		case StatusMissing:
			missing++
		}
	}
	if installed != 2 {
		t.Errorf("expected 2 installed, got %d", installed)
	}
	if missing != 5 {
		t.Errorf("expected 5 missing, got %d", missing)
	}
}

func TestInstall_Brew_Success(t *testing.T) {
	runner := &mockRunner{
		lookPath: map[string]error{"brew": nil},
		runs: map[string]runResult{
			"brew install neovim": {output: []byte("installed")},
		},
	}
	inst := NewInstaller(runner, "darwin")

	err := inst.Install(App{Name: "Neovim", BinName: "nvim", BrewPkg: "neovim"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestInstall_Brew_NoPkg(t *testing.T) {
	runner := &mockRunner{
		lookPath: map[string]error{},
		runs:     map[string]runResult{},
	}
	inst := NewInstaller(runner, "darwin")

	err := inst.Install(App{Name: "Test", BinName: "test", BrewPkg: ""})

	if err == nil {
		t.Error("expected error for empty BrewPkg")
	}
	if !strings.Contains(err.Error(), "no homebrew package") {
		t.Errorf("expected 'no homebrew package' error, got: %v", err)
	}
}

func TestInstall_Brew_NotFound(t *testing.T) {
	runner := &mockRunner{
		lookPath: map[string]error{},
		runs:     map[string]runResult{},
	}
	inst := NewInstaller(runner, "darwin")

	err := inst.Install(App{Name: "Neovim", BinName: "nvim", BrewPkg: "neovim"})

	if err == nil {
		t.Error("expected error when brew not found")
	}
	if !strings.Contains(err.Error(), "homebrew not found") {
		t.Errorf("expected 'homebrew not found' error, got: %v", err)
	}
}

func TestInstall_Brew_Failure(t *testing.T) {
	runner := &mockRunner{
		lookPath: map[string]error{"brew": nil},
		runs: map[string]runResult{
			"brew install neovim": {output: []byte("Error: something failed"), err: fmt.Errorf("exit status 1")},
		},
	}
	inst := NewInstaller(runner, "darwin")

	err := inst.Install(App{Name: "Neovim", BinName: "nvim", BrewPkg: "neovim"})

	if err == nil {
		t.Error("expected error on brew failure")
	}
	if !strings.Contains(err.Error(), "brew install neovim") {
		t.Errorf("expected 'brew install neovim' in error, got: %v", err)
	}
}

func TestInstall_Apt_Success(t *testing.T) {
	runner := &mockRunner{
		lookPath: map[string]error{"apt-get": nil},
		runs: map[string]runResult{
			"sudo apt-get install -y neovim": {output: []byte("installed")},
		},
	}
	inst := NewInstaller(runner, "linux")

	err := inst.Install(App{Name: "Neovim", BinName: "nvim", AptPkg: "neovim"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestInstall_Apt_NoPkg(t *testing.T) {
	runner := &mockRunner{
		lookPath: map[string]error{},
		runs:     map[string]runResult{},
	}
	inst := NewInstaller(runner, "linux")

	err := inst.Install(App{Name: "Oh My Posh", BinName: "oh-my-posh", AptPkg: ""})

	if err == nil {
		t.Error("expected error for empty AptPkg")
	}
	if !strings.Contains(err.Error(), "no apt package") {
		t.Errorf("expected 'no apt package' error, got: %v", err)
	}
}

func TestInstall_Apt_NotFound(t *testing.T) {
	runner := &mockRunner{
		lookPath: map[string]error{},
		runs:     map[string]runResult{},
	}
	inst := NewInstaller(runner, "linux")

	err := inst.Install(App{Name: "Neovim", BinName: "nvim", AptPkg: "neovim"})

	if err == nil {
		t.Error("expected error when apt-get not found")
	}
	if !strings.Contains(err.Error(), "apt-get not found") {
		t.Errorf("expected 'apt-get not found' error, got: %v", err)
	}
}

func TestInstall_UnsupportedOS(t *testing.T) {
	runner := &mockRunner{
		lookPath: map[string]error{},
		runs:     map[string]runResult{},
	}
	inst := NewInstaller(runner, "windows")

	err := inst.Install(App{Name: "Neovim", BinName: "nvim", BrewPkg: "neovim", AptPkg: "neovim"})

	if err == nil {
		t.Error("expected error for unsupported OS")
	}
	if !strings.Contains(err.Error(), "unsupported OS") {
		t.Errorf("expected 'unsupported OS' error, got: %v", err)
	}
}

func TestNewInstaller_DefaultOS(t *testing.T) {
	runner := &mockRunner{
		lookPath: map[string]error{},
		runs:     map[string]runResult{},
	}
	inst := NewInstaller(runner, "")

	if inst.os == "" {
		t.Error("expected os to be set to runtime.GOOS")
	}
}
