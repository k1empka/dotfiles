# dotfiles

Cross-platform dotfiles managed with [chezmoi](https://www.chezmoi.io/), themed with [oh-my-posh](https://ohmyposh.dev/).

## Supported Platforms

| OS | Shell | Prompt |
|----|-------|--------|
| macOS | ZSH | oh-my-posh |
| Linux | ZSH / Bash | oh-my-posh |
| Windows | PowerShell | oh-my-posh |

## Quick Install

### macOS / Linux

```sh
curl -fsLS https://raw.githubusercontent.com/k1empka/dotfiles/refs/heads/main/scripts/bootstrap.sh | sh
```

### Windows (PowerShell)
.
```powershell
winget install twpayne.chezmoi
winget install JanDeDobbeleer.OhMyPosh
chezmoi init --apply k1empka
```

## What's Included

### Shell Configurations
- **ZSH** — history, completion, key bindings, aliases, plugin support (macOS/Linux)
- **Bash** — fallback config with matching aliases (macOS/Linux)
- **PowerShell** — PSReadLine, Terminal-Icons, matching aliases (Windows)

### Prompt Themes
All platforms use oh-my-posh with configurable themes:
- **Rose Pine** (default)
- **Catppuccin Mocha**
- **Tokyo Night**

Each theme includes segments for: OS icon, path, git status, Node.js/Go/Python versions, execution time, exit code, and clock.

### TUI Manager (WIP)
A lazygit-style terminal UI for managing dotfiles, built with Go + [bubbletea](https://github.com/charmbracelet/bubbletea).

```sh
make build
./bin/dotfiles-tui
```

Tabs: Overview, Shell, Git, Themes, Neovim, Alacritty, Chezmoi, Install, VS Code

The **Install** tab detects which applications are installed and supports one-click installation via Homebrew (macOS) or apt (Linux). The **VS Code** tab provides a split-pane browser for viewing managed VS Code configuration files.

## Configuration

During `chezmoi init`, you'll be prompted for:

| Variable | Description |
|----------|-------------|
| `name` | Your full name (for git config) |
| `email` | Your email address |
| `theme` | Prompt theme: `rose-pine`, `catppuccin`, or `tokyo-night` |

### Changing Theme

```sh
chezmoi edit-config   # change theme in [data] section
chezmoi apply         # apply the new theme
```

Then restart your shell.

### Local Overrides

Machine-specific config that won't be tracked by git:

| Shell | File |
|-------|------|
| ZSH | `~/.zshrc.local` |
| Bash | `~/.bashrc.local` |
| PowerShell | `~/.powershell.local.ps1` |

## Repository Structure

```
dotfiles/
├── cmd/dotfiles-tui/       # Go TUI entry point
├── internal/
│   ├── tui/                # TUI components (bubbletea)
│   ├── chezmoi/            # chezmoi CLI wrapper
│   └── installer/          # application installer (brew/apt)
├── home/                   # chezmoi source directory
│   ├── .chezmoi.toml.tmpl  # chezmoi config template
│   ├── .chezmoiignore      # OS-conditional ignore rules
│   ├── dot_zshrc.tmpl      # ZSH config
│   ├── dot_bashrc.tmpl     # Bash config
│   ├── dot_config/
│   │   └── oh-my-posh/
│   │       └── themes/     # oh-my-posh theme files
│   ├── private_Library/    # VS Code config (macOS only)
│   │   └── .../Code/User/
│   │       ├── settings.json
│   │       └── extensions.txt
│   └── Documents/PowerShell/
│       └── Microsoft.PowerShell_profile.ps1.tmpl
├── scripts/
│   └── bootstrap.sh        # one-liner bootstrap
├── go.mod
├── Makefile
└── README.md
```

## Prerequisites

- [chezmoi](https://www.chezmoi.io/install/)
- [oh-my-posh](https://ohmyposh.dev/docs/installation/linux)
- A [Nerd Font](https://www.nerdfonts.com/) (e.g., JetBrainsMono Nerd Font)
- Go 1.24+ (only for building the TUI)

## Development

```sh
make build    # build TUI binary
make test     # run tests
make lint     # run go vet
```
