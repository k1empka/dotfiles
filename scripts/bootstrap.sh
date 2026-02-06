#!/bin/sh
set -e

# Bootstrap script for setting up dotfiles on a new machine.
# Usage: curl -fsLS https://raw.githubusercontent.com/k1empka/dotfiles/main/scripts/bootstrap.sh | sh

REPO="k1empka/dotfiles"
NERD_FONT="JetBrainsMono"
NERD_FONT_VERSION="3.3.0"

echo "==> Bootstrapping dotfiles..."

# --- Helpers ---

require_brew() {
    if ! command -v brew >/dev/null 2>&1; then
        echo "    Homebrew not found. Install it first: https://brew.sh"
        exit 1
    fi
}

install_if_missing() {
    bin="$1"
    label="$2"
    shift 2

    if command -v "$bin" >/dev/null 2>&1; then
        echo "==> ${label} already installed"
        return
    fi

    echo "==> Installing ${label}..."
    "$@"
}

# --- Nerd Font ---

install_nerd_font() {
    FONT_DIR=""
    case "$(uname -s)" in
        Darwin) FONT_DIR="$HOME/Library/Fonts" ;;
        Linux)  FONT_DIR="$HOME/.local/share/fonts" ;;
        *)
            echo "    Unsupported OS for font install."
            return
            ;;
    esac

    if ls "${FONT_DIR}/"*"${NERD_FONT}"* >/dev/null 2>&1; then
        echo "==> ${NERD_FONT} Nerd Font already installed"
        return
    fi

    echo "==> Installing ${NERD_FONT} Nerd Font..."
    mkdir -p "$FONT_DIR"
    ZIPFILE="/tmp/${NERD_FONT}.zip"
    curl -fsSL -o "$ZIPFILE" \
        "https://github.com/ryanoasis/nerd-fonts/releases/download/v${NERD_FONT_VERSION}/${NERD_FONT}.zip"
    unzip -oq "$ZIPFILE" -d "$FONT_DIR"
    rm -f "$ZIPFILE"

    # Rebuild font cache on Linux.
    if [ "$(uname -s)" = "Linux" ] && command -v fc-cache >/dev/null 2>&1; then
        fc-cache -f "$FONT_DIR"
    fi

    echo "    ${NERD_FONT} Nerd Font installed to ${FONT_DIR}"
}

# --- Git ---

install_git() {
    install_if_missing git "Git" _install_git_impl
}

_install_git_impl() {
    case "$(uname -s)" in
        Darwin)
            require_brew
            brew install git
            ;;
        Linux)
            if command -v apt-get >/dev/null 2>&1; then
                sudo apt-get install -y git
            elif command -v pacman >/dev/null 2>&1; then
                sudo pacman -S --noconfirm git
            fi
            ;;
    esac
}

# --- ZSH ---

install_zsh() {
    install_if_missing zsh "ZSH" _install_zsh_impl
}

_install_zsh_impl() {
    case "$(uname -s)" in
        Darwin)
            require_brew
            brew install zsh
            ;;
        Linux)
            if command -v apt-get >/dev/null 2>&1; then
                sudo apt-get install -y zsh
            elif command -v pacman >/dev/null 2>&1; then
                sudo pacman -S --noconfirm zsh
            fi
            ;;
    esac
}

# --- ZSH plugins ---

install_zsh_plugins() {
    echo "==> Installing ZSH plugins..."
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

# --- Oh My Posh ---

install_oh_my_posh() {
    install_if_missing oh-my-posh "Oh My Posh" _install_oh_my_posh_impl
}

_install_oh_my_posh_impl() {
    case "$(uname -s)" in
        Darwin)
            require_brew
            brew install jandedobbeleer/oh-my-posh/oh-my-posh
            ;;
        Linux)
            curl -s https://ohmyposh.dev/install.sh | bash -s
            ;;
        *)
            echo "    Unsupported OS. On Windows, use: winget install JanDeDobbeleer.OhMyPosh"
            exit 1
            ;;
    esac
}

# --- Neovim ---

install_neovim() {
    install_if_missing nvim "Neovim" _install_neovim_impl
}

_install_neovim_impl() {
    case "$(uname -s)" in
        Darwin)
            require_brew
            brew install neovim
            ;;
        Linux)
            if command -v brew >/dev/null 2>&1; then
                brew install neovim
            else
                echo "    Installing neovim AppImage..."
                curl -fsSL -o /tmp/nvim.appimage \
                    https://github.com/neovim/neovim/releases/latest/download/nvim.appimage
                chmod u+x /tmp/nvim.appimage
                sudo mkdir -p /usr/local/bin
                sudo mv /tmp/nvim.appimage /usr/local/bin/nvim
            fi
            ;;
    esac
}

# --- Alacritty ---

install_alacritty() {
    install_if_missing alacritty "Alacritty" _install_alacritty_impl
}

_install_alacritty_impl() {
    case "$(uname -s)" in
        Darwin)
            require_brew
            brew install alacritty
            ;;
        Linux)
            if command -v apt-get >/dev/null 2>&1; then
                sudo apt-get install -y alacritty
            elif command -v pacman >/dev/null 2>&1; then
                sudo pacman -S --noconfirm alacritty
            fi
            ;;
    esac
}

# --- Run all ---

install_nerd_font
install_git
install_zsh
install_zsh_plugins
install_oh_my_posh
install_neovim
install_alacritty

# Install chezmoi and apply dotfiles.
sh -c "$(curl -fsLS get.chezmoi.io)" -- init --apply "$REPO"

echo "==> Done! Restart your shell to see changes."
