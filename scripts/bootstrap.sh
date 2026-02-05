#!/bin/sh
set -e

# Bootstrap script for setting up dotfiles on a new machine.
# Usage: curl -fsLS https://raw.githubusercontent.com/k1empka/dotfiles/main/scripts/bootstrap.sh | sh

REPO="k1empka/dotfiles"

echo "==> Installing chezmoi and applying dotfiles..."

# Detect OS and install oh-my-posh
install_oh_my_posh() {
    if command -v oh-my-posh >/dev/null 2>&1; then
        echo "==> oh-my-posh already installed"
        return
    fi

    echo "==> Installing oh-my-posh..."
    case "$(uname -s)" in
        Darwin)
            if command -v brew >/dev/null 2>&1; then
                brew install jandedobbeleer/oh-my-posh/oh-my-posh
            else
                echo "    Homebrew not found. Install it first: https://brew.sh"
                exit 1
            fi
            ;;
        Linux)
            curl -s https://ohmyposh.dev/install.sh | bash -s
            ;;
        *)
            echo "    Unsupported OS for this script. On Windows, use:"
            echo "    winget install JanDeDobbeleer.OhMyPosh"
            exit 1
            ;;
    esac
}

# Install ZSH plugins (macOS/Linux)
install_zsh_plugins() {
    case "$(uname -s)" in
        Darwin)
            if command -v brew >/dev/null 2>&1; then
                brew install zsh-autosuggestions zsh-syntax-highlighting
            fi
            ;;
        Linux)
            if command -v apt-get >/dev/null 2>&1; then
                sudo apt-get install -y zsh-autosuggestions zsh-syntax-highlighting 2>/dev/null || true
            elif command -v pacman >/dev/null 2>&1; then
                sudo pacman -S --noconfirm zsh-autosuggestions zsh-syntax-highlighting 2>/dev/null || true
            fi
            ;;
    esac
}

# Install neovim
install_neovim() {
    if command -v nvim >/dev/null 2>&1; then
        echo "==> neovim already installed"
        return
    fi

    echo "==> Installing neovim..."
    case "$(uname -s)" in
        Darwin)
            if command -v brew >/dev/null 2>&1; then
                brew install neovim
            else
                echo "    Homebrew not found. Install it first: https://brew.sh"
                exit 1
            fi
            ;;
        Linux)
            if command -v brew >/dev/null 2>&1; then
                brew install neovim
            else
                echo "    Installing neovim AppImage..."
                curl -LO https://github.com/neovim/neovim/releases/latest/download/nvim.appimage
                chmod u+x nvim.appimage
                sudo mkdir -p /usr/local/bin
                sudo mv nvim.appimage /usr/local/bin/nvim
            fi
            ;;
        *)
            echo "    Unsupported OS for this script."
            exit 1
            ;;
    esac
}

install_oh_my_posh
install_zsh_plugins
install_neovim

# Install chezmoi and apply dotfiles
sh -c "$(curl -fsLS get.chezmoi.io)" -- init --apply "$REPO"

echo "==> Done! Restart your shell to see changes."
