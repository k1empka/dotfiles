package chezmoi

import (
	"os"
	"path/filepath"
	"testing"
)

func TestThemeColors_ExtractsUniqueColors(t *testing.T) {
	dir := t.TempDir()
	themesDir := filepath.Join(dir, "dot_config", "oh-my-posh", "themes")
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		t.Fatal(err)
	}

	themeJSON := `{
		"blocks": [{
			"segments": [
				{"foreground": "#e0def4", "background": "#31748f"},
				{"foreground": "#e0def4", "background": "#c4a7e7"},
				{"foreground": "#191724", "background": "#ebbcba"}
			]
		}]
	}`
	themeFile := filepath.Join(themesDir, "test-theme.omp.json")
	if err := os.WriteFile(themeFile, []byte(themeJSON), 0644); err != nil {
		t.Fatal(err)
	}

	colors, err := extractThemeColors(themeFile)
	if err != nil {
		t.Fatalf("extractThemeColors: %v", err)
	}

	// Should have 5 unique colors.
	if len(colors) != 5 {
		t.Errorf("expected 5 unique colors, got %d: %v", len(colors), colors)
	}

	// Verify specific colors present.
	colorSet := make(map[string]bool)
	for _, c := range colors {
		colorSet[c] = true
	}
	for _, want := range []string{"#e0def4", "#31748f", "#c4a7e7", "#191724", "#ebbcba"} {
		if !colorSet[want] {
			t.Errorf("expected color %s in result", want)
		}
	}
}

func TestThemeColors_EmptyTheme(t *testing.T) {
	dir := t.TempDir()
	themeFile := filepath.Join(dir, "empty.omp.json")
	if err := os.WriteFile(themeFile, []byte(`{"blocks":[]}`), 0644); err != nil {
		t.Fatal(err)
	}

	colors, err := extractThemeColors(themeFile)
	if err != nil {
		t.Fatalf("extractThemeColors: %v", err)
	}
	if len(colors) != 0 {
		t.Errorf("expected 0 colors, got %d", len(colors))
	}
}
