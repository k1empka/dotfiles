package installer

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// Status represents the installation state of an app.
type Status int

const (
	StatusUnknown Status = iota
	StatusInstalled
	StatusMissing
	StatusInstalling
	StatusError
)

// AppStatus holds an app and its current state.
type AppStatus struct {
	App     App
	Status  Status
	Version string
	Err     error
}

// CommandRunner abstracts command execution for testability.
type CommandRunner interface {
	LookPath(file string) (string, error)
	Run(name string, args ...string) ([]byte, error)
}

// ExecRunner is the real implementation using os/exec.
type ExecRunner struct{}

// LookPath finds a binary in PATH.
func (ExecRunner) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

// Run executes a command and returns combined output.
func (ExecRunner) Run(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

// Installer checks and installs applications.
type Installer struct {
	runner CommandRunner
	os     string
}

// NewInstaller creates an Installer.
// The osName parameter should be "darwin" or "linux"; empty defaults to runtime.GOOS.
func NewInstaller(runner CommandRunner, osName string) *Installer {
	if osName == "" {
		osName = runtime.GOOS
	}
	return &Installer{runner: runner, os: osName}
}

// CheckAll checks the install status of all apps.
func (inst *Installer) CheckAll() []AppStatus {
	apps := Apps()
	results := make([]AppStatus, len(apps))
	for i, app := range apps {
		results[i] = inst.Check(app)
	}
	return results
}

// validateName rejects names containing shell metacharacters or path separators.
func validateName(name string) error {
	for _, c := range name {
		switch {
		case c == '/' || c == '\\' || c == ';' || c == '&' || c == '|' || c == '`' || c == '$' || c == '(' || c == ')':
			return fmt.Errorf("invalid character %q in name %q", string(c), name)
		}
	}
	return nil
}

// Check determines whether a single app is installed.
func (inst *Installer) Check(app App) AppStatus {
	if err := validateName(app.BinName); err != nil {
		return AppStatus{App: app, Status: StatusError, Err: err}
	}
	_, err := inst.runner.LookPath(app.BinName)
	if err != nil {
		return AppStatus{App: app, Status: StatusMissing}
	}

	version := ""
	out, err := inst.runner.Run(app.BinName, "--version")
	if err == nil {
		lines := strings.SplitN(strings.TrimSpace(string(out)), "\n", 2)
		if len(lines) > 0 {
			version = lines[0]
		}
	}

	return AppStatus{App: app, Status: StatusInstalled, Version: version}
}

// Install installs an app using the appropriate package manager.
func (inst *Installer) Install(app App) error {
	if err := validateName(app.BinName); err != nil {
		return fmt.Errorf("invalid app: %w", err)
	}
	switch inst.os {
	case "darwin":
		return inst.installBrew(app)
	case "linux":
		return inst.installApt(app)
	default:
		return fmt.Errorf("unsupported OS: %s", inst.os)
	}
}

// brewArgs returns the arguments for a brew install command.
// Packages prefixed with "--cask " are installed as casks.
func brewArgs(pkg string) []string {
	if name, ok := strings.CutPrefix(pkg, "--cask "); ok {
		return []string{"install", "--cask", name}
	}
	return []string{"install", pkg}
}

func (inst *Installer) installBrew(app App) error {
	if app.BrewPkg == "" {
		return fmt.Errorf("no homebrew package defined for %s", app.Name)
	}
	_, err := inst.runner.LookPath("brew")
	if err != nil {
		return fmt.Errorf("homebrew not found in PATH: %w", err)
	}
	args := brewArgs(app.BrewPkg)
	out, err := inst.runner.Run("brew", args...)
	if err != nil {
		return fmt.Errorf("brew %s: %s: %w", strings.Join(args, " "), strings.TrimSpace(string(out)), err)
	}
	return nil
}

func (inst *Installer) installApt(app App) error {
	if app.AptPkg == "" {
		return fmt.Errorf("no apt package defined for %s", app.Name)
	}
	_, err := inst.runner.LookPath("apt-get")
	if err != nil {
		return fmt.Errorf("apt-get not found in PATH: %w", err)
	}
	out, err := inst.runner.Run("sudo", "apt-get", "install", "-y", app.AptPkg)
	if err != nil {
		return fmt.Errorf("apt-get install %s: %s: %w", app.AptPkg, strings.TrimSpace(string(out)), err)
	}
	return nil
}
