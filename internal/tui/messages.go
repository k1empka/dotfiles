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

// ConfigUpdatedMsg signals that chezmoi config was written and applied.
type ConfigUpdatedMsg struct {
	Err error
}

// ErrorMsg is a generic error message.
type ErrorMsg struct {
	Err error
}
