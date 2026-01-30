#!/bin/bash

# Common utilities for dotfiles installation scripts
# Source this file in other scripts: source "$(dirname "$0")/common.sh"

# --- Colors ---
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# --- Output Helpers ---
info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[OK]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# --- OS Detection ---
detect_os() {
    OS="$(uname -s)"
    case "$OS" in
        Linux*)     MACHINE=Linux;;
        Darwin*)    MACHINE=Mac;;
        *)          MACHINE="UNKNOWN:${OS}";;
    esac
    export MACHINE

    # Detect Linux distribution
    if [ "$MACHINE" = "Linux" ] && [ -f /etc/os-release ]; then
        . /etc/os-release
        export DISTRO_ID="$ID"
    else
        export DISTRO_ID=""
    fi
}

# --- Utility Functions ---
command_exists() {
    command -v "$1" &> /dev/null
}

ensure_dir() {
    if [ ! -d "$1" ]; then
        mkdir -p "$1"
        info "Created directory: $1"
    fi
}

# --- Symlink Management ---
# Creates a symlink with backup of existing file
# Usage: safe_symlink <source> <target> [dry_run]
safe_symlink() {
    local source="$1"
    local target="$2"
    local dry_run="${3:-false}"
    local backup_dir="$HOME/dotfiles_old"

    # Validate source exists
    if [ ! -e "$source" ]; then
        error "Source does not exist: $source"
        return 1
    fi

    # If target is already a correct symlink, skip
    if [ -L "$target" ] && [ "$(readlink "$target")" = "$source" ]; then
        success "Already linked: $target -> $source"
        return 0
    fi

    if [ "$dry_run" = "true" ]; then
        info "[DRY-RUN] Would link: $target -> $source"
        return 0
    fi

    # Backup existing file/directory if it exists and is not a symlink
    if [ -e "$target" ] || [ -L "$target" ]; then
        ensure_dir "$backup_dir"
        local backup_name
        backup_name="$(basename "$target").$(date +%Y%m%d%H%M%S)"
        mv "$target" "$backup_dir/$backup_name"
        warn "Backed up existing: $target -> $backup_dir/$backup_name"
    fi

    # Ensure parent directory exists
    ensure_dir "$(dirname "$target")"

    # Create symlink
    ln -s "$source" "$target"
    success "Linked: $target -> $source"
}

# --- Script Directory ---
# Get the directory where the dotfiles repo is located
get_dotfiles_dir() {
    local script_dir
    script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    # scripts/ is one level down from dotfiles root
    dirname "$script_dir"
}

# Initialize OS detection when sourced
detect_os
