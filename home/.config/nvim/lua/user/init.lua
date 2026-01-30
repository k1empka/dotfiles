---@type LazySpec
return {
  "rose-pine/neovim",
  config = function(_, opts)
    require("rose-pine").setup(opts)
    vim.cmd("colorscheme rose-pine")
  end,
}
