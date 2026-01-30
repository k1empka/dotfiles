# Bash Configuration

# --- Early Exit for Non-Interactive Shells ---
case $- in
    *i*) ;;
      *) return;;
esac

# --- History Settings ---
HISTSIZE=10000
HISTFILESIZE=20000
HISTCONTROL=ignoreboth:erasedups
shopt -s histappend

# --- Shell Options ---
shopt -s checkwinsize
shopt -s cdspell
shopt -s dirspell 2>/dev/null
shopt -s autocd 2>/dev/null
shopt -s globstar 2>/dev/null

# --- Colors ---
export CLICOLOR=1
export LSCOLORS=GxFxCxDxBxegedabagaced

# Enable color support
if [ -x /usr/bin/dircolors ]; then
    test -r ~/.dircolors && eval "$(dircolors -b ~/.dircolors)" || eval "$(dircolors -b)"
fi

# --- Prompt ---
# Simple colored prompt showing user@host:directory
if [ "$TERM" != "dumb" ]; then
    PS1='\[\033[01;32m\]\u@\h\[\033[00m\]:\[\033[01;34m\]\w\[\033[00m\]\$ '
else
    PS1='\u@\h:\w\$ '
fi

# --- PATH ---
# Add local bin directories
[ -d "$HOME/.local/bin" ] && export PATH="$HOME/.local/bin:$PATH"
[ -d "$HOME/bin" ] && export PATH="$HOME/bin:$PATH"

# Neovim on Linux (if installed from tarball)
[ -d "/opt/nvim-linux-x86_64/bin" ] && export PATH="$PATH:/opt/nvim-linux-x86_64/bin"

# --- Default Editor ---
if command -v nvim &> /dev/null; then
    export EDITOR=nvim
    export VISUAL=nvim
elif command -v vim &> /dev/null; then
    export EDITOR=vim
    export VISUAL=vim
fi

# --- Aliases ---

# Navigation
alias ..='cd ..'
alias ...='cd ../..'
alias ....='cd ../../..'

# List files
alias ls='ls --color=auto 2>/dev/null || ls -G'
alias ll='ls -alF'
alias la='ls -A'
alias l='ls -CF'
alias lt='ls -altr'

# Safety nets
alias rm='rm -i'
alias cp='cp -i'
alias mv='mv -i'

# Grep with color
alias grep='grep --color=auto'
alias fgrep='fgrep --color=auto'
alias egrep='egrep --color=auto'

# Git shortcuts
alias g='git'
alias gs='git status'
alias ga='git add'
alias gc='git commit'
alias gp='git push'
alias gl='git pull'
alias gd='git diff'
alias gco='git checkout'
alias gb='git branch'
alias glog='git log --oneline --graph --decorate -20'

# Neovim/Vim
if command -v nvim &> /dev/null; then
    alias vim='nvim'
    alias vi='nvim'
fi

# Clear screen
alias cls='clear'

# Show disk usage
alias df='df -h'
alias du='du -h'

# Process management
alias psg='ps aux | grep -v grep | grep -i'

# --- Functions ---

# Create directory and cd into it
mkcd() {
    mkdir -p "$1" && cd "$1"
}

# Extract various archive formats
extract() {
    if [ -f "$1" ]; then
        case "$1" in
            *.tar.bz2)   tar xjf "$1"     ;;
            *.tar.gz)    tar xzf "$1"     ;;
            *.tar.xz)    tar xJf "$1"     ;;
            *.bz2)       bunzip2 "$1"     ;;
            *.rar)       unrar x "$1"     ;;
            *.gz)        gunzip "$1"      ;;
            *.tar)       tar xf "$1"      ;;
            *.tbz2)      tar xjf "$1"     ;;
            *.tgz)       tar xzf "$1"     ;;
            *.zip)       unzip "$1"       ;;
            *.Z)         uncompress "$1"  ;;
            *.7z)        7z x "$1"        ;;
            *)           echo "'$1' cannot be extracted via extract()" ;;
        esac
    else
        echo "'$1' is not a valid file"
    fi
}

# Find file by name
ff() {
    find . -type f -iname "*$1*"
}

# Find directory by name
fd() {
    find . -type d -iname "*$1*"
}

# --- Completions ---
# Enable bash completion if available
if ! shopt -oq posix; then
    if [ -f /usr/share/bash-completion/bash_completion ]; then
        . /usr/share/bash-completion/bash_completion
    elif [ -f /etc/bash_completion ]; then
        . /etc/bash_completion
    elif [ -f /opt/homebrew/etc/profile.d/bash_completion.sh ]; then
        . /opt/homebrew/etc/profile.d/bash_completion.sh
    fi
fi

# --- Local Customizations ---
# Source local bashrc if it exists (for machine-specific settings)
[ -f ~/.bashrc.local ] && source ~/.bashrc.local
