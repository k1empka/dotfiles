package chezmoi

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ThemeColors extracts unique foreground and background colors from an OMP theme file.
func (c *Client) ThemeColors(themeName string) ([]string, error) {
	srcPath, err := c.SourcePath()
	if err != nil {
		return nil, err
	}
	themeFile := filepath.Join(srcPath, "dot_config", "oh-my-posh", "themes", themeName+".omp.json")
	return extractThemeColors(themeFile)
}

// extractThemeColors reads an OMP theme JSON file and returns unique colors.
func extractThemeColors(path string) ([]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read theme: %w", err)
	}

	var theme struct {
		Blocks []struct {
			Segments []struct {
				Foreground string `json:"foreground"`
				Background string `json:"background"`
			} `json:"segments"`
		} `json:"blocks"`
	}
	if err := json.Unmarshal(content, &theme); err != nil {
		return nil, fmt.Errorf("failed to parse theme: %w", err)
	}

	seen := make(map[string]bool)
	var colors []string
	for _, block := range theme.Blocks {
		for _, seg := range block.Segments {
			for _, c := range []string{seg.Foreground, seg.Background} {
				if c != "" && !seen[c] {
					seen[c] = true
					colors = append(colors, c)
				}
			}
		}
	}
	return colors, nil
}
