package chezmoi

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListSourceFiles(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "subdir")
	if err := os.MkdirAll(sub, 0755); err != nil {
		t.Fatal(err)
	}

	// Create test files.
	for _, name := range []string{"a.txt", "b.txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(sub, "c.txt"), []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	// Walk the directory.
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(dir, path)
		files = append(files, rel)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 3 {
		t.Errorf("expected 3 files, got %d: %v", len(files), files)
	}
}

func TestReadFile(t *testing.T) {
	dir := t.TempDir()
	testFile := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(testFile, []byte("hello world"), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "hello world" {
		t.Errorf("expected 'hello world', got %q", content)
	}
}

func TestThemeNames_FromDir(t *testing.T) {
	dir := t.TempDir()

	// Create mock theme files.
	for _, name := range []string{"rose-pine.omp.json", "catppuccin.omp.json", "readme.txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("{}"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	var names []string
	for _, e := range entries {
		if name, found := cutSuffix(e.Name(), ".omp.json"); found {
			names = append(names, name)
		}
	}

	if len(names) != 2 {
		t.Errorf("expected 2 themes, got %d: %v", len(names), names)
	}
}

func TestAlacrittyConfig_ReadFile(t *testing.T) {
	dir := t.TempDir()
	alDir := filepath.Join(dir, "dot_config", "alacritty")
	if err := os.MkdirAll(alDir, 0755); err != nil {
		t.Fatal(err)
	}

	content := `import = ["~/.config/alacritty/themes/rose-pine.toml"]`
	testFile := filepath.Join(alDir, "alacritty.toml.tmpl")
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != content {
		t.Errorf("expected %q, got %q", content, got)
	}
}

func cutSuffix(s, suffix string) (string, bool) {
	if len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix {
		return s[:len(s)-len(suffix)], true
	}
	return s, false
}
