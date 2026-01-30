#!/bin/bash

# Symlink creation script
# Creates symlinks from home directory to dotfiles

set -e

_LINK_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$_LINK_DIR/common.sh"

DRY_RUN="${DRY_RUN:-false}"

# --- Configuration ---

# Get dotfiles directory
DOTFILES_DIR="$(get_dotfiles_dir)"
HOME_DIR="$DOTFILES_DIR/home"

# Files to symlink directly in home directory
HOME_FILES=(
    ".bashrc"
    ".zshrc"
    ".vimrc"
    ".gitconfig"
    ".tmux.conf"
)

# --- Symlink Functions ---

link_home_files() {
    info "Linking home directory files..."

    for file in "${HOME_FILES[@]}"; do
        local source="$HOME_DIR/$file"
        local target="$HOME/$file"

        if [ -f "$source" ]; then
            safe_symlink "$source" "$target" "$DRY_RUN"
        else
            warn "Source file not found, skipping: $source"
        fi
    done
}

link_neovim_config() {
    info "Linking Neovim user configuration..."

    local source="$HOME_DIR/.config/nvim/lua/user"
    local target="$HOME/.config/nvim/lua/user"

    if [ ! -d "$source" ]; then
        warn "Neovim user config not found, skipping: $source"
        return 0
    fi

    # Ensure nvim config directory exists
    ensure_dir "$HOME/.config/nvim/lua"

    safe_symlink "$source" "$target" "$DRY_RUN"
}

link_alacritty_config() {
    info "Linking Alacritty configuration..."

    local alacritty_source="$HOME_DIR/.config/alacritty"
    local alacritty_target="$HOME/.config/alacritty"

    if [ ! -d "$alacritty_source" ]; then
        warn "Alacritty config not found, skipping: $alacritty_source"
        return 0
    fi

    ensure_dir "$alacritty_target"

    # Link alacritty.toml
    if [ -f "$alacritty_source/alacritty.toml" ]; then
        safe_symlink "$alacritty_source/alacritty.toml" "$alacritty_target/alacritty.toml" "$DRY_RUN"
    fi

    # Link downloads directory (contains themes)
    if [ -d "$alacritty_source/downloads" ]; then
        safe_symlink "$alacritty_source/downloads" "$alacritty_target/downloads" "$DRY_RUN"
    fi
}

# --- Main ---

show_summary() {
    echo ""
    info "=== Symlink Summary ==="
    echo ""
    info "Home files:"
    for file in "${HOME_FILES[@]}"; do
        echo "  $HOME/$file -> $HOME_DIR/$file"
    done
    echo ""
    info "Neovim config:"
    echo "  $HOME/.config/nvim/lua/user -> $HOME_DIR/.config/nvim/lua/user"
    echo ""
    info "Alacritty config:"
    echo "  $HOME/.config/alacritty/alacritty.toml -> $HOME_DIR/.config/alacritty/alacritty.toml"
    echo "  $HOME/.config/alacritty/downloads -> $HOME_DIR/.config/alacritty/downloads"
    echo ""
}

create_links() {
    info "Creating symlinks..."
    info "Dotfiles directory: $DOTFILES_DIR"
    info "Dry run: $DRY_RUN"
    echo ""

    link_home_files
    link_neovim_config
    link_alacritty_config

    echo ""
    success "Symlink creation complete!"

    if [ "$DRY_RUN" = "true" ]; then
        warn "This was a dry run. No changes were made."
        warn "Run without --dry-run to apply changes."
    fi
}

# Run if executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    # Check for dry-run flag
    if [[ "$1" == "--dry-run" ]] || [[ "$1" == "-n" ]]; then
        DRY_RUN=true
    fi

    if [[ "$1" == "--summary" ]]; then
        show_summary
    else
        create_links
    fi
fi
