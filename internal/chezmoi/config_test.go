package chezmoi

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteConfig_UpdatesExistingKeys(t *testing.T) {
	dir := t.TempDir()
	dataFile := filepath.Join(dir, ".chezmoidata.toml")
	initial := `name = "Alice"
email = "alice@example.com"
theme = "rose-pine"
`
	if err := os.WriteFile(dataFile, []byte(initial), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a mock client that returns our temp dir as source path.
	c := &Client{bin: "echo"} // bin doesn't matter, we override SourcePath behavior
	// Since WriteConfig calls SourcePath which shells out, we test the TOML logic directly.
	err := writeConfigToFile(dataFile, map[string]string{
		"name":  "Bob",
		"theme": "catppuccin",
	})
	if err != nil {
		t.Fatalf("writeConfigToFile: %v", err)
	}

	content, err := os.ReadFile(dataFile)
	if err != nil {
		t.Fatal(err)
	}
	_ = c // suppress unused

	got := string(content)
	tests := []struct {
		name string
		want string
	}{
		{"name updated", `name = "Bob"`},
		{"email preserved", `email = "alice@example.com"`},
		{"theme updated", `theme = "catppuccin"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !contains(got, tt.want) {
				t.Errorf("expected %q in:\n%s", tt.want, got)
			}
		})
	}
}

func TestWriteConfig_AppendsNewKeys(t *testing.T) {
	dir := t.TempDir()
	dataFile := filepath.Join(dir, ".chezmoidata.toml")
	initial := `name = "Alice"
`
	if err := os.WriteFile(dataFile, []byte(initial), 0644); err != nil {
		t.Fatal(err)
	}

	err := writeConfigToFile(dataFile, map[string]string{
		"email": "bob@example.com",
	})
	if err != nil {
		t.Fatalf("writeConfigToFile: %v", err)
	}

	content, err := os.ReadFile(dataFile)
	if err != nil {
		t.Fatal(err)
	}

	got := string(content)
	if !contains(got, `email = "bob@example.com"`) {
		t.Errorf("expected appended email key in:\n%s", got)
	}
	if !contains(got, `name = "Alice"`) {
		t.Errorf("expected name preserved in:\n%s", got)
	}
}

func TestWriteConfig_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	dataFile := filepath.Join(dir, ".chezmoidata.toml")
	if err := os.WriteFile(dataFile, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	err := writeConfigToFile(dataFile, map[string]string{
		"name": "Charlie",
	})
	if err != nil {
		t.Fatalf("writeConfigToFile: %v", err)
	}

	content, err := os.ReadFile(dataFile)
	if err != nil {
		t.Fatal(err)
	}

	got := string(content)
	if !contains(got, `name = "Charlie"`) {
		t.Errorf("expected name in:\n%s", got)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
