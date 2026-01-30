# Dotfiles

Personal development environment configuration with a unified Rose Pine theme.

## What's Included

| Tool | Configuration | Theme |
|------|--------------|-------|
| Bash | `.bashrc` | - |
| Zsh | `.zshrc` + Oh My Zsh | - |
| Git | `.gitconfig` | - |
| Vim | `.vimrc` | slate |
| Neovim | AstroNvim + user config | Rose Pine |
| Alacritty | `alacritty.toml` | Rose Pine |
| Tmux | `.tmux.conf` | Rose Pine |

## Quick Start

```bash
# Clone the repository
git clone https://github.com/yourusername/dotfiles.git ~/dotfiles
cd ~/dotfiles

# Install everything
./install.sh
```

## Installation Options

The installer supports selective installation via flags:

```bash
# Show help
./install.sh --help

# Install everything (default)
./install.sh
./install.sh -a

# Install only dependencies
./install.sh -d

# Install only frameworks (Oh My Zsh, AstroNvim)
./install.sh -f

# Create symlinks only
./install.sh -l

# Preview changes without making them
./install.sh --dry-run
./install.sh -l -n    # Preview symlinks only

# Combine flags
./install.sh -d -f    # Dependencies and frameworks, skip symlinks
```

## Dependencies Installed

The script installs these tools (if not already present):

- **git** - Version control
- **neovim** - Text editor
- **vim** - Classic editor
- **zsh** - Z shell
- **tmux** - Terminal multiplexer
- **alacritty** - GPU-accelerated terminal
- **lazygit** - Terminal UI for git

### Platform Support

| Platform | Package Manager |
|----------|----------------|
| macOS | Homebrew |
| Ubuntu/Debian | apt + PPA |
| Arch Linux | pacman |
| Fedora/CentOS | dnf + copr |

## Directory Structure

```
dotfiles/
├── install.sh              # Main installer with flags
├── scripts/
│   ├── common.sh          # Shared utilities
│   ├── deps.sh            # Dependency installation
│   ├── frameworks.sh      # Oh My Zsh, AstroNvim
│   └── link.sh            # Symlink management
└── home/
    ├── .bashrc            # Bash configuration
    ├── .zshrc             # Zsh configuration
    ├── .vimrc             # Vim configuration
    ├── .gitconfig         # Git configuration
    ├── .tmux.conf         # Tmux configuration
    └── .config/
        ├── alacritty/     # Alacritty terminal config
        └── nvim/          # Neovim user config
```

## Post-Installation

After installation, complete these steps:

1. **Restart your terminal** or run:
   ```bash
   source ~/.zshrc  # or ~/.bashrc
   ```

2. **Open Neovim** to install plugins:
   ```bash
   nvim
   ```

3. **Configure Git** with your identity:
   ```bash
   git config --global user.name "Your Name"
   git config --global user.email "your@email.com"
   ```

## Customization

### Local Overrides

Create these files for machine-specific settings (not tracked by git):

- `~/.bashrc.local` - Local bash customizations
- `~/.zshrc.local` - Local zsh customizations (add to .zshrc if needed)

### Modifying Configs

Edit files in `~/dotfiles/home/` - they're symlinked to your home directory.

## Backup

Existing configurations are backed up to `~/dotfiles_old/` before being replaced. Each backup includes a timestamp.

## Troubleshooting

### Symlinks not working

```bash
# Check if symlinks are correct
ls -la ~/.bashrc ~/.zshrc ~/.gitconfig

# Re-run symlink creation
./install.sh -l
```

### Neovim plugins not loading

```bash
# Open Neovim and run
:Lazy sync
```

### Oh My Zsh not loading

Ensure `.zshrc` sources Oh My Zsh:
```bash
grep "oh-my-zsh.sh" ~/.zshrc
```

## Git Aliases

The `.gitconfig` includes these useful aliases:

| Alias | Command | Description |
|-------|---------|-------------|
| `st` | `status` | Show status |
| `co` | `checkout` | Switch branches |
| `br` | `branch` | List branches |
| `ci` | `commit` | Commit changes |
| `lg` | `log --oneline --graph` | Pretty log |
| `df` | `diff` | Show changes |
| `unstage` | `reset HEAD --` | Unstage files |
| `undo` | `reset --soft HEAD~1` | Undo last commit |

Run `git aliases` to see all available aliases.

## License

MIT
