#!/bin/bash

# Dependency installation script
# Installs: git, neovim, alacritty, zsh, lazygit, tmux

set -e

_DEPS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$_DEPS_DIR/common.sh"

# --- Package Lists ---
PACKAGES="git neovim alacritty zsh lazygit tmux"

# --- Installation Functions ---

install_macos() {
    info "Installing dependencies for macOS..."

    # Check for Homebrew
    if ! command_exists brew; then
        info "Homebrew not found. Installing Homebrew..."
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    fi

    info "Updating Homebrew..."
    brew update

    for pkg in $PACKAGES; do
        if brew list "$pkg" &>/dev/null; then
            success "$pkg is already installed"
        else
            info "Installing $pkg..."
            brew install "$pkg"
        fi
    done
}

install_ubuntu() {
    info "Installing dependencies for Ubuntu..."

    sudo apt-get update

    # Install base packages
    sudo apt-get install -y software-properties-common curl

    # Add lazygit PPA (Ubuntu only)
    if ! grep -q "lazygit-team" /etc/apt/sources.list.d/* 2>/dev/null; then
        sudo add-apt-repository ppa:lazygit-team/release -y
        sudo apt-get update
    fi

    # Install packages available via apt
    local apt_packages="git zsh tmux lazygit"
    for pkg in $apt_packages; do
        if dpkg -l "$pkg" &>/dev/null 2>&1; then
            success "$pkg is already installed"
        else
            info "Installing $pkg..."
            sudo apt-get install -y "$pkg"
        fi
    done

    # Alacritty may not be in default repos
    install_alacritty_apt

    # Install Neovim from official release
    install_neovim_linux
}

install_debian() {
    info "Installing dependencies for Debian..."

    sudo apt-get update

    # Install base packages
    sudo apt-get install -y curl

    # Install packages available via apt
    local apt_packages="git zsh tmux"
    for pkg in $apt_packages; do
        if dpkg -l "$pkg" &>/dev/null 2>&1; then
            success "$pkg is already installed"
        else
            info "Installing $pkg..."
            sudo apt-get install -y "$pkg"
        fi
    done

    # Install lazygit from GitHub releases (PPAs don't work on Debian)
    install_lazygit_github

    # Alacritty may not be in default repos
    install_alacritty_apt

    # Install Neovim from official release
    install_neovim_linux
}

install_lazygit_github() {
    if command_exists lazygit; then
        success "lazygit is already installed"
        return 0
    fi

    info "Installing lazygit from GitHub releases..."

    local lazygit_version
    lazygit_version=$(curl -s "https://api.github.com/repos/jesseduffield/lazygit/releases/latest" | grep -Po '"tag_name": "v\K[^"]*')

    if [ -z "$lazygit_version" ]; then
        warn "Could not determine lazygit version, using fallback"
        lazygit_version="0.44.1"
    fi

    local lazygit_tar="/tmp/lazygit.tar.gz"
    curl -Lo "$lazygit_tar" "https://github.com/jesseduffield/lazygit/releases/latest/download/lazygit_${lazygit_version}_Linux_x86_64.tar.gz"
    sudo tar -C /usr/local/bin -xzf "$lazygit_tar" lazygit
    rm "$lazygit_tar"

    success "lazygit installed to /usr/local/bin"
}

install_alacritty_apt() {
    if command_exists alacritty; then
        success "alacritty is already installed"
        return 0
    fi

    # Try to install from apt (may not be available on all versions)
    if apt-cache show alacritty &>/dev/null; then
        info "Installing alacritty..."
        sudo apt-get install -y alacritty
    else
        warn "alacritty not available in apt repos. Install manually from: https://github.com/alacritty/alacritty"
    fi
}

install_arch() {
    info "Installing dependencies for Arch Linux..."

    sudo pacman -Syu --noconfirm

    local pacman_packages="git zsh alacritty lazygit tmux"
    for pkg in $pacman_packages; do
        if pacman -Qi "$pkg" &>/dev/null; then
            success "$pkg is already installed"
        else
            info "Installing $pkg..."
            sudo pacman -S --noconfirm "$pkg"
        fi
    done

    # Install Neovim from official release
    install_neovim_linux
}

install_fedora() {
    info "Installing dependencies for Fedora/CentOS..."

    sudo dnf update -y

    # Add lazygit copr repo
    if ! dnf copr list | grep -q "atim/lazygit"; then
        sudo dnf copr enable atim/lazygit -y
    fi

    local dnf_packages="git zsh alacritty lazygit tmux"
    for pkg in $dnf_packages; do
        if rpm -q "$pkg" &>/dev/null; then
            success "$pkg is already installed"
        else
            info "Installing $pkg..."
            sudo dnf install -y "$pkg"
        fi
    done

    # Install Neovim from official release
    install_neovim_linux
}

install_neovim_linux() {
    if command_exists nvim; then
        success "Neovim is already installed"
        return 0
    fi

    info "Installing Neovim from official release..."

    local nvim_url="https://github.com/neovim/neovim/releases/latest/download/nvim-linux-x86_64.tar.gz"
    local nvim_tar="/tmp/nvim-linux-x86_64.tar.gz"
    local nvim_dir="/opt/nvim-linux-x86_64"

    curl -Lo "$nvim_tar" "$nvim_url"
    sudo rm -rf /opt/nvim /opt/nvim-linux-x86_64
    sudo tar -C /opt -xzf "$nvim_tar"
    rm "$nvim_tar"

    # Add to PATH in profile if not already there
    local profile_file="$HOME/.profile"
    local path_entry='export PATH="$PATH:/opt/nvim-linux-x86_64/bin"'

    if ! grep -q '/opt/nvim-linux-x86_64/bin' "$profile_file" 2>/dev/null; then
        echo "$path_entry" >> "$profile_file"
        info "Added Neovim to PATH in $profile_file"
    fi

    success "Neovim installed to $nvim_dir"
}

# --- Main ---

check_deps() {
    info "Checking installed dependencies..."
    local missing=0

    for pkg in $PACKAGES; do
        if command_exists "$pkg"; then
            success "$pkg is installed"
        else
            warn "$pkg is NOT installed"
            missing=1
        fi
    done

    return $missing
}

install_deps() {
    info "Detected OS: $MACHINE"

    if [ "$MACHINE" = "Mac" ]; then
        install_macos
    elif [ "$MACHINE" = "Linux" ]; then
        case "$DISTRO_ID" in
            ubuntu|linuxmint|pop)
                install_ubuntu
                ;;
            debian)
                install_debian
                ;;
            arch|manjaro)
                install_arch
                ;;
            fedora|centos|rhel)
                install_fedora
                ;;
            *)
                error "Unsupported Linux distribution: $DISTRO_ID"
                exit 1
                ;;
        esac
    else
        error "Unsupported OS: $MACHINE"
        exit 1
    fi

    success "All dependencies installed!"
}

# Run if executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    install_deps
fi
