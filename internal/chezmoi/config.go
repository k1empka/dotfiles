package chezmoi

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// WriteConfig updates chezmoi template data values in the .chezmoidata.toml file.
func (c *Client) WriteConfig(values map[string]string) error {
	srcPath, err := c.SourcePath()
	if err != nil {
		return fmt.Errorf("failed to get source path: %w", err)
	}
	return writeConfigToFile(filepath.Join(srcPath, ".chezmoidata.toml"), values)
}

// writeConfigToFile performs targeted key=value replacement in a TOML file.
func writeConfigToFile(path string, values map[string]string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", path, err)
	}

	updated := string(content)
	for key, val := range values {
		pattern := regexp.MustCompile(`(?m)^(\s*` + regexp.QuoteMeta(key) + `\s*=\s*).*$`)
		replacement := fmt.Sprintf(`${1}%q`, val)
		if pattern.MatchString(updated) {
			updated = pattern.ReplaceAllString(updated, replacement)
		} else {
			updated = strings.TrimRight(updated, "\n") + "\n" + key + " = " + fmt.Sprintf("%q", val) + "\n"
		}
	}

	if err := os.WriteFile(path, []byte(updated), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}
	return nil
}

// WriteConfigAndApply writes config values and runs chezmoi apply.
func (c *Client) WriteConfigAndApply(values map[string]string) error {
	if err := c.WriteConfig(values); err != nil {
		return err
	}
	return c.Apply()
}
