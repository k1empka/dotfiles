package installer

// App describes an installable application.
type App struct {
	Name    string // Display name, e.g. "Neovim"
	BinName string // Binary to check in PATH, e.g. "nvim"
	BrewPkg string // Homebrew formula, e.g. "neovim"
	AptPkg  string // apt package name (empty if not available via apt)
}

// Apps returns the list of supported applications.
func Apps() []App {
	return []App{
		{Name: "Neovim", BinName: "nvim", BrewPkg: "neovim", AptPkg: "neovim"},
		{Name: "Alacritty", BinName: "alacritty", BrewPkg: "alacritty", AptPkg: "alacritty"},
		{Name: "Git", BinName: "git", BrewPkg: "git", AptPkg: "git"},
		{Name: "ZSH", BinName: "zsh", BrewPkg: "zsh", AptPkg: "zsh"},
		{Name: "Bash", BinName: "bash", BrewPkg: "bash", AptPkg: "bash"},
		{Name: "Oh My Posh", BinName: "oh-my-posh", BrewPkg: "jandedobbeleer/oh-my-posh/oh-my-posh", AptPkg: ""},
		{Name: "Chezmoi", BinName: "chezmoi", BrewPkg: "chezmoi", AptPkg: ""},
		{Name: "VS Code", BinName: "code", BrewPkg: "--cask visual-studio-code", AptPkg: "code"},
	}
}
