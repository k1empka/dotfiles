package chezmoi

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ReadSourceFile reads a file from the chezmoi source directory.
func (c *Client) ReadSourceFile(relPath string) (string, error) {
	srcPath, err := c.SourcePath()
	if err != nil {
		return "", err
	}
	content, err := os.ReadFile(filepath.Join(srcPath, relPath))
	if err != nil {
		return "", fmt.Errorf("failed to read source file %s: %w", relPath, err)
	}
	return string(content), nil
}

// ListSourceFiles returns relative paths of files under a directory in the source tree.
func (c *Client) ListSourceFiles(dir string) ([]string, error) {
	srcPath, err := c.SourcePath()
	if err != nil {
		return nil, err
	}

	root := filepath.Join(srcPath, dir)
	var files []string
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files = append(files, rel)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list source files in %s: %w", dir, err)
	}
	return files, nil
}

// ShellConfigs returns the content of each shell config template.
// Returns [ZSH, Bash, PowerShell] contents.
func (c *Client) ShellConfigs() ([3]string, error) {
	var configs [3]string
	shellFiles := [3]string{
		"dot_zshrc.tmpl",
		"dot_bashrc.tmpl",
		"Documents/PowerShell/Microsoft.PowerShell_profile.ps1.tmpl",
	}
	for i, f := range shellFiles {
		content, err := c.ReadSourceFile(f)
		if err != nil {
			// Not all shells exist on all platforms.
			configs[i] = fmt.Sprintf("(file not found: %s)", f)
			continue
		}
		configs[i] = content
	}
	return configs, nil
}

// GitConfig returns the gitconfig template content.
func (c *Client) GitConfig() (string, error) {
	return c.ReadSourceFile("dot_gitconfig.tmpl")
}

// NvimFiles returns the list of neovim config files.
func (c *Client) NvimFiles() ([]string, error) {
	return c.ListSourceFiles("dot_config/nvim")
}

// NvimFileContent returns the content of a specific neovim config file.
func (c *Client) NvimFileContent(relPath string) (string, error) {
	return c.ReadSourceFile(filepath.Join("dot_config/nvim", relPath))
}

// AlacrittyConfig returns the alacritty config template content.
func (c *Client) AlacrittyConfig() (string, error) {
	return c.ReadSourceFile("dot_config/alacritty/alacritty.toml.tmpl")
}

// ThemeNames returns available OMP theme names by scanning the themes directory.
func (c *Client) ThemeNames() ([]string, error) {
	srcPath, err := c.SourcePath()
	if err != nil {
		return nil, err
	}

	themesDir := filepath.Join(srcPath, "dot_config", "oh-my-posh", "themes")
	entries, err := os.ReadDir(themesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read themes directory: %w", err)
	}

	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".omp.json") {
			names = append(names, strings.TrimSuffix(name, ".omp.json"))
		}
	}
	return names, nil
}
