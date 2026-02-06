package tui

// ConfigLoadedMsg carries the initial chezmoi configuration data.
type ConfigLoadedMsg struct {
	Name    string
	Email   string
	Theme   string
	OS      string
	Version string
	Err     error
}

// ErrorMsg is a generic error message.
type ErrorMsg struct {
	Err error
}
