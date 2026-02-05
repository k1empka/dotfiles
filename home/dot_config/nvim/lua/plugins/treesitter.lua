---@type LazySpec
return {
  "nvim-treesitter/nvim-treesitter",
  opts = {
    ensure_installed = {
      "lua",
      "vim",
      "vimdoc",
      "go",
      "bash",
      "json",
      "yaml",
      "markdown",
    },
  },
}
