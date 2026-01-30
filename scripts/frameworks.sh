#!/bin/bash

# Framework installation script
# Installs: Oh My Zsh, AstroNvim

set -e

_FRAMEWORKS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$_FRAMEWORKS_DIR/common.sh"

# --- Zsh ---

install_zsh() {
    if command_exists zsh; then
        success "zsh is already installed"
        return 0
    fi

    info "Installing zsh..."

    if [ "$MACHINE" = "Mac" ]; then
        if command_exists brew; then
            brew install zsh
        else
            error "Homebrew not found. Please install Homebrew first."
            return 1
        fi
    elif [ "$MACHINE" = "Linux" ]; then
        case "$DISTRO_ID" in
            ubuntu|debian|linuxmint|pop)
                sudo apt-get update
                sudo apt-get install -y zsh
                ;;
            arch|manjaro)
                sudo pacman -S --noconfirm zsh
                ;;
            fedora|centos|rhel)
                sudo dnf install -y zsh
                ;;
            *)
                error "Unsupported distribution: $DISTRO_ID. Please install zsh manually."
                return 1
                ;;
        esac
    else
        error "Unsupported OS: $MACHINE. Please install zsh manually."
        return 1
    fi

    success "zsh installed!"
}

# --- Oh My Zsh ---

install_ohmyzsh() {
    if [ -d "$HOME/.oh-my-zsh" ]; then
        success "Oh My Zsh is already installed"
        return 0
    fi

    # Ensure zsh is installed first
    install_zsh

    info "Installing Oh My Zsh..."

    # Install Oh My Zsh unattended (won't change shell automatically)
    RUNZSH=no CHSH=no sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"

    success "Oh My Zsh installed!"
}

# --- AstroNvim ---

install_astronvim() {
    local nvim_config="$HOME/.config/nvim"

    if [ -d "$nvim_config" ]; then
        # Check if it's AstroNvim
        if [ -f "$nvim_config/lua/astronvim/init.lua" ] || grep -q "AstroNvim" "$nvim_config/init.lua" 2>/dev/null; then
            success "AstroNvim is already installed"
            return 0
        fi

        # Check if it's the AstroNvim template
        if [ -d "$nvim_config/.git" ]; then
            local remote
            remote=$(git -C "$nvim_config" remote get-url origin 2>/dev/null || echo "")
            if [[ "$remote" == *"AstroNvim"* ]]; then
                success "AstroNvim template is already installed"
                return 0
            fi
        fi

        warn "Existing Neovim config found at $nvim_config"
        warn "Backing up to $HOME/dotfiles_old/"
        ensure_dir "$HOME/dotfiles_old"
        mv "$nvim_config" "$HOME/dotfiles_old/nvim.$(date +%Y%m%d%H%M%S)"
    fi

    info "Installing AstroNvim template..."
    git clone --depth 1 https://github.com/AstroNvim/template "$nvim_config"

    # Remove the template's git history so user can make it their own
    rm -rf "$nvim_config/.git"

    success "AstroNvim template installed!"
    info "Run 'nvim' to complete setup and install plugins"
}

# --- Main ---

install_frameworks() {
    info "Installing frameworks..."

    install_ohmyzsh
    install_astronvim

    success "All frameworks installed!"
}

# Run if executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    install_frameworks
fi
