package chezmoi

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// Client wraps chezmoi CLI commands.
type Client struct {
	bin string
}

// NewClient creates a Client, locating the chezmoi binary.
func NewClient() (*Client, error) {
	bin, err := exec.LookPath("chezmoi")
	if err != nil {
		return nil, fmt.Errorf("chezmoi not found in PATH: %w", err)
	}
	return &Client{bin: bin}, nil
}

// ConfigData holds the user-configured chezmoi template data.
type ConfigData struct {
	Name  string
	Email string
	Theme string
	OS    string
}

// Data returns the chezmoi template data via `chezmoi data --format=json`.
func (c *Client) Data() (*ConfigData, error) {
	out, err := exec.Command(c.bin, "data", "--format=json").CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("chezmoi data: %w", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(out, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse chezmoi data: %w", err)
	}

	data := &ConfigData{}

	// Extract user data fields.
	if name, ok := raw["name"].(string); ok {
		data.Name = name
	}
	if email, ok := raw["email"].(string); ok {
		data.Email = email
	}
	if theme, ok := raw["theme"].(string); ok {
		data.Theme = theme
	}

	// Extract OS from chezmoi metadata.
	if chezmoiData, ok := raw["chezmoi"].(map[string]any); ok {
		if osVal, ok := chezmoiData["os"].(string); ok {
			data.OS = osVal
		}
	}

	return data, nil
}

// Version returns the chezmoi version string.
func (c *Client) Version() (string, error) {
	out, err := exec.Command(c.bin, "--version").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("chezmoi version: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// Status returns the output of `chezmoi status`.
func (c *Client) Status() (string, error) {
	out, err := exec.Command(c.bin, "status").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("chezmoi status: %w", err)
	}
	return string(out), nil
}

// Apply runs `chezmoi apply`.
func (c *Client) Apply() error {
	cmd := exec.Command(c.bin, "apply")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("chezmoi apply: %w", err)
	}
	return nil
}

// Diff returns the output of `chezmoi diff`.
func (c *Client) Diff() (string, error) {
	out, err := exec.Command(c.bin, "diff").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("chezmoi diff: %w", err)
	}
	return string(out), nil
}

// SourcePath returns the chezmoi source directory path.
func (c *Client) SourcePath() (string, error) {
	out, err := exec.Command(c.bin, "source-path").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("chezmoi source-path: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
