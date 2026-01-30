package chezmoi

import (
	"fmt"
	"os/exec"
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
