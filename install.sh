#!/bin/bash

# Dotfiles Installation Script
# A modular installer for personal development environment

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/scripts/common.sh"

# --- Default Options ---
INSTALL_DEPS=false
INSTALL_FRAMEWORKS=false
INSTALL_LINKS=false
DRY_RUN=false
SHOW_HELP=false

# --- Usage ---
usage() {
    cat << EOF
Dotfiles Installation Script

Usage: $(basename "$0") [OPTIONS]

Options:
  -a, --all          Install everything (default if no options given)
  -d, --deps         Install dependencies only (git, nvim, zsh, etc.)
  -f, --frameworks   Install frameworks only (Oh My Zsh, AstroNvim)
  -l, --link         Create symlinks only
  -n, --dry-run      Show what would be done without making changes
  -h, --help         Show this help message

Examples:
  $(basename "$0")              # Install everything
  $(basename "$0") -a           # Install everything (explicit)
  $(basename "$0") -d           # Install dependencies only
  $(basename "$0") -l -n        # Preview symlink changes
  $(basename "$0") -d -f        # Install deps and frameworks, skip symlinks

EOF
}

# --- Argument Parsing ---
parse_args() {
    # If no arguments, install everything
    if [ $# -eq 0 ]; then
        INSTALL_DEPS=true
        INSTALL_FRAMEWORKS=true
        INSTALL_LINKS=true
        return
    fi

    while [[ $# -gt 0 ]]; do
        case $1 in
            -a|--all)
                INSTALL_DEPS=true
                INSTALL_FRAMEWORKS=true
                INSTALL_LINKS=true
                shift
                ;;
            -d|--deps)
                INSTALL_DEPS=true
                shift
                ;;
            -f|--frameworks)
                INSTALL_FRAMEWORKS=true
                shift
                ;;
            -l|--link)
                INSTALL_LINKS=true
                shift
                ;;
            -n|--dry-run)
                DRY_RUN=true
                shift
                ;;
            -h|--help)
                SHOW_HELP=true
                shift
                ;;
            *)
                error "Unknown option: $1"
                usage
                exit 1
                ;;
        esac
    done
}

# --- Main Installation ---
main() {
    parse_args "$@"

    if [ "$SHOW_HELP" = true ]; then
        usage
        exit 0
    fi

    echo ""
    info "=========================================="
    info "       Dotfiles Installation Script       "
    info "=========================================="
    echo ""
    info "Detected OS: $MACHINE"
    [ -n "$DISTRO_ID" ] && info "Distribution: $DISTRO_ID"
    echo ""

    if [ "$DRY_RUN" = true ]; then
        warn "DRY RUN MODE - No changes will be made"
        echo ""
    fi

    # Install dependencies
    if [ "$INSTALL_DEPS" = true ]; then
        info ">>> Installing Dependencies..."
        if [ "$DRY_RUN" = true ]; then
            info "[DRY-RUN] Would install: git, neovim, alacritty, zsh, lazygit, tmux"
        else
            source "$SCRIPT_DIR/scripts/deps.sh"
            install_deps
        fi
        echo ""
    fi

    # Install frameworks
    if [ "$INSTALL_FRAMEWORKS" = true ]; then
        info ">>> Installing Frameworks..."
        if [ "$DRY_RUN" = true ]; then
            info "[DRY-RUN] Would install: Oh My Zsh, AstroNvim"
        else
            source "$SCRIPT_DIR/scripts/frameworks.sh"
            install_frameworks
        fi
        echo ""
    fi

    # Create symlinks
    if [ "$INSTALL_LINKS" = true ]; then
        info ">>> Creating Symlinks..."
        export DRY_RUN
        source "$SCRIPT_DIR/scripts/link.sh"
        create_links
        echo ""
    fi

    # Summary
    echo ""
    info "=========================================="
    if [ "$DRY_RUN" = true ]; then
        warn "DRY RUN COMPLETE - No changes were made"
    else
        success "Installation Complete!"
        echo ""
        info "Next steps:"
        info "  1. Restart your terminal or run: source ~/.zshrc"
        info "  2. Open Neovim to install plugins: nvim"
        info "  3. Update your git config: git config --global user.name 'Your Name'"
        info "     and: git config --global user.email 'your@email.com'"
    fi
    info "=========================================="
}

main "$@"
