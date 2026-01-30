" Vim Configuration
" Basic settings for vanilla Vim (Neovim users should use AstroNvim)

" --- General Settings ---

" Disable compatibility mode
set nocompatible

" Enable file type detection
filetype plugin indent on

" Enable syntax highlighting
syntax enable

" Use UTF-8 encoding
set encoding=utf-8
set fileencoding=utf-8

" --- Display ---

" Show line numbers
set number
set relativenumber

" Highlight current line
set cursorline

" Show matching brackets
set showmatch

" Always show status line
set laststatus=2

" Show command in bottom bar
set showcmd

" Show current mode
set showmode

" Enable wildmenu for command completion
set wildmenu
set wildmode=longest:full,full

" Disable error bells
set noerrorbells
set visualbell
set t_vb=

" --- Search ---

" Highlight search results
set hlsearch

" Incremental search
set incsearch

" Case insensitive search unless uppercase used
set ignorecase
set smartcase

" --- Indentation ---

" Use spaces instead of tabs
set expandtab

" Tab width
set tabstop=4
set shiftwidth=4
set softtabstop=4

" Smart indentation
set smartindent
set autoindent

" --- Editing ---

" Enable backspace in insert mode
set backspace=indent,eol,start

" Keep some lines visible when scrolling
set scrolloff=8
set sidescrolloff=8

" Don't wrap lines
set nowrap

" Enable hidden buffers
set hidden

" Faster updates
set updatetime=300

" Don't create backup files
set nobackup
set nowritebackup
set noswapfile

" Persistent undo
set undofile
set undodir=~/.vim/undodir

" --- Splits ---

" Open new splits below and to the right
set splitbelow
set splitright

" --- Key Mappings ---

" Set leader key to space
let mapleader = " "

" Quick save
nnoremap <leader>w :w<CR>

" Quick quit
nnoremap <leader>q :q<CR>

" Clear search highlighting
nnoremap <leader>h :nohlsearch<CR>

" Navigate splits with Ctrl + hjkl
nnoremap <C-h> <C-w>h
nnoremap <C-j> <C-w>j
nnoremap <C-k> <C-w>k
nnoremap <C-l> <C-w>l

" Move lines up/down in visual mode
vnoremap J :m '>+1<CR>gv=gv
vnoremap K :m '<-2<CR>gv=gv

" Keep cursor centered when jumping
nnoremap n nzzzv
nnoremap N Nzzzv

" Better indenting in visual mode
vnoremap < <gv
vnoremap > >gv

" --- Colors ---

" Enable true colors if supported
if has('termguicolors')
    set termguicolors
endif

" Dark background
set background=dark

" Try to use a nice colorscheme
try
    colorscheme slate
catch
    " Fall back to default
endtry
