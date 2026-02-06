package panel

// SetSizeMsg is sent to panels when available content area changes.
type SetSizeMsg struct {
	Width  int
	Height int
}

// ConfigDataMsg carries chezmoi configuration data to panels.
type ConfigDataMsg struct {
	Name    string
	Email   string
	Theme   string
	OS      string
	Version string
}

// ConfigUpdatedMsg signals that a config update completed.
type ConfigUpdatedMsg struct {
	Err error
}

// ShellContentMsg carries shell config file content.
type ShellContentMsg struct {
	Index   int // 0=ZSH, 1=Bash, 2=PowerShell
	Content string
}

// GitContentMsg carries gitconfig content.
type GitContentMsg struct {
	Content string
}

// AlacrittyContentMsg carries alacritty config content.
type AlacrittyContentMsg struct {
	Content string
}

// NvimFilesMsg carries the list of neovim config files.
type NvimFilesMsg struct {
	Files []string
}

// NvimContentMsg carries neovim file content.
type NvimContentMsg struct {
	Path    string
	Content string
}

// ChezmoiOutputMsg carries chezmoi command output.
type ChezmoiOutputMsg struct {
	Mode   string // "status", "diff", "apply"
	Output string
	Err    error
}

// ThemeColorsMsg carries extracted theme colors.
type ThemeColorsMsg struct {
	Theme  string
	Colors []string
}

// ThemeApplyMsg signals the app to apply a theme.
type ThemeApplyMsg struct {
	Theme string
}

// RunChezmoiMsg signals the app to run a chezmoi command.
type RunChezmoiMsg struct {
	Mode string // "status", "diff", "apply"
}
