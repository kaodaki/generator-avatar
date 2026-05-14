package avatar

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	cfg := DefaultConfig()
	cfg.OutputPath = t.TempDir()

	a, err := New(cfg)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("New() returned nil")
	}
}

func TestNew_NilConfig(t *testing.T) {
	a, err := New(nil)
	if err != nil {
		t.Fatalf("New(nil) unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("New(nil) returned nil")
	}
}

func TestGenerate(t *testing.T) {
	cfg := DefaultConfig()
	cfg.OutputPath = t.TempDir()
	cfg.OutputName = "lib_test"
	cfg.EmailTruncateLength = 100 // Большое значение чтобы не сокращать

	a, err := New(cfg)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	result, err := a.Generate("library@example.com")
	if err != nil {
		t.Fatalf("Generate() unexpected error: %v", err)
	}

	if result.FilePath == "" {
		t.Error("FilePath is empty")
	}
	if result.BackgroundColor == "" {
		t.Error("BackgroundColor is empty")
	}
	if result.Text != "l" {
		t.Errorf("Text=%s, want 'l'", result.Text)
	}

	// Проверяем что файл существует
	if _, err := os.Stat(result.FilePath); os.IsNotExist(err) {
		t.Errorf("file does not exist: %s", result.FilePath)
	}
}

func TestGenerateWithOptions(t *testing.T) {
	cfg := DefaultConfig()
	cfg.OutputPath = t.TempDir()

	a, _ := New(cfg)

	customPath := t.TempDir()
	result, err := a.GenerateWithOptions("options@test.com", customPath, "custom_name")
	if err != nil {
		t.Fatalf("GenerateWithOptions() unexpected error: %v", err)
	}

	expectedPath := filepath.Join(customPath, "custom_name.png")
	if result.FilePath != expectedPath {
		t.Errorf("FilePath=%s, want %s", result.FilePath, expectedPath)
	}

	if _, err := os.Stat(result.FilePath); os.IsNotExist(err) {
		t.Errorf("file does not exist: %s", result.FilePath)
	}
}

func TestSetOutputPath(t *testing.T) {
	cfg := DefaultConfig()
	cfg.OutputPath = t.TempDir()

	a, _ := New(cfg)

	newPath := filepath.Join(t.TempDir(), "subdir")
	a.SetOutputPath(newPath)
	a.SetOutputName("path_change_test")

	result, err := a.Generate("path@test.com")
	if err != nil {
		t.Fatalf("Generate() unexpected error: %v", err)
	}

	if filepath.Dir(result.FilePath) != newPath {
		t.Errorf("FilePath dir=%s, want %s", filepath.Dir(result.FilePath), newPath)
	}
}

func TestSetOutputName(t *testing.T) {
	cfg := DefaultConfig()
	cfg.OutputPath = t.TempDir()

	a, _ := New(cfg)
	a.SetOutputName("renamed_avatar")

	result, err := a.Generate("rename@test.com")
	if err != nil {
		t.Fatalf("Generate() unexpected error: %v", err)
	}

	if filepath.Base(result.FilePath) != "renamed_avatar.png" {
		t.Errorf("FilePath basename=%s, want renamed_avatar.png", filepath.Base(result.FilePath))
	}
}

func TestUpdateConfig(t *testing.T) {
	cfg := DefaultConfig()
	cfg.OutputPath = t.TempDir()

	a, _ := New(cfg)

	newCfg := DefaultConfig()
	newCfg.OutputPath = t.TempDir()
	newCfg.Width = 600
	newCfg.Height = 600
	newCfg.FontSize = 80

	err := a.UpdateConfig(newCfg)
	if err != nil {
		t.Fatalf("UpdateConfig() unexpected error: %v", err)
	}

	a.SetOutputName("updated_config")
	result, err := a.Generate("update@test.com")
	if err != nil {
		t.Fatalf("Generate() unexpected error: %v", err)
	}

	if result.FilePath == "" {
		t.Error("FilePath is empty after config update")
	}
}

func TestPaletteCount(t *testing.T) {
	cfg := DefaultConfig()
	cfg.OutputPath = t.TempDir()

	a, _ := New(cfg)

	count := a.PaletteCount()
	if count != 12 {
		t.Errorf("PaletteCount()=%d, want 12", count)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg == nil {
		t.Fatal("DefaultConfig() returned nil")
	}
	if cfg.Width != 300 {
		t.Errorf("Width=%d, want 300", cfg.Width)
	}
	if cfg.Height != 300 {
		t.Errorf("Height=%d, want 300", cfg.Height)
	}
	if cfg.FontSize != 45 {
		t.Errorf("FontSize=%f, want 45", cfg.FontSize)
	}
}

// Test helper functions

func TestWithWidth(t *testing.T) {
	cfg := DefaultConfig()
	WithWidth(cfg, 500)
	if cfg.Width != 500 {
		t.Errorf("Width=%d, want 500", cfg.Width)
	}
}

func TestWithHeight(t *testing.T) {
	cfg := DefaultConfig()
	WithHeight(cfg, 500)
	if cfg.Height != 500 {
		t.Errorf("Height=%d, want 500", cfg.Height)
	}
}

func TestWithFontSize(t *testing.T) {
	cfg := DefaultConfig()
	WithFontSize(cfg, 60)
	if cfg.FontSize != 60 {
		t.Errorf("FontSize=%f, want 60", cfg.FontSize)
	}
}

func TestWithOutputPath(t *testing.T) {
	cfg := DefaultConfig()
	WithOutputPath(cfg, "/custom/path")
	if cfg.OutputPath != "/custom/path" {
		t.Errorf("OutputPath=%s, want /custom/path", cfg.OutputPath)
	}
}

func TestWithOutputName(t *testing.T) {
	cfg := DefaultConfig()
	WithOutputName(cfg, "custom")
	if cfg.OutputName != "custom" {
		t.Errorf("OutputName=%s, want custom", cfg.OutputName)
	}
}

func TestWithEmailTruncateLength(t *testing.T) {
	cfg := DefaultConfig()
	WithEmailTruncateLength(cfg, 5)
	if cfg.EmailTruncateLength != 5 {
		t.Errorf("EmailTruncateLength=%d, want 5", cfg.EmailTruncateLength)
	}
}

// Test ConfigOption builder

func TestNewConfig_WithOptions(t *testing.T) {
	cfg := NewConfig(
		ConfigWithWidth(400),
		ConfigWithHeight(400),
		ConfigWithFontSize(55),
		ConfigWithEmailTruncateLength(4),
	)

	if cfg.Width != 400 {
		t.Errorf("Width=%d, want 400", cfg.Width)
	}
	if cfg.Height != 400 {
		t.Errorf("Height=%d, want 400", cfg.Height)
	}
	if cfg.FontSize != 55 {
		t.Errorf("FontSize=%f, want 55", cfg.FontSize)
	}
	if cfg.EmailTruncateLength != 4 {
		t.Errorf("EmailTruncateLength=%d, want 4", cfg.EmailTruncateLength)
	}
}

func TestNewConfig_Default(t *testing.T) {
	cfg := NewConfig()

	if cfg.Width != 300 {
		t.Errorf("Width=%d, want 300", cfg.Width)
	}
	if cfg.Height != 300 {
		t.Errorf("Height=%d, want 300", cfg.Height)
	}
}

// Test utility functions

func TestPickRandomColor(t *testing.T) {
	c := PickRandomColor()
	if c == nil {
		t.Fatal("PickRandomColor() returned nil")
	}
	if c.Hex == "" {
		t.Error("PickRandomColor() returned color with empty Hex")
	}
}

func TestLoadPaletteFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	colorsFile := filepath.Join(tmpDir, "test_colors.json")

	jsonData := `[
		{"name": "Red", "hex": "#FF0000", "r": "255", "g": "0", "b": "0"}
	]`
	os.WriteFile(colorsFile, []byte(jsonData), 0644)

	palette, err := LoadPaletteFromFile(colorsFile)
	if err != nil {
		t.Fatalf("LoadPaletteFromFile() unexpected error: %v", err)
	}
	if palette.Count() != 1 {
		t.Errorf("Count=%d, want 1", palette.Count())
	}
}

func TestLoadPaletteFromFile_NotFound(t *testing.T) {
	_, err := LoadPaletteFromFile("/nonexistent/file.json")
	if err == nil {
		t.Error("LoadPaletteFromFile() expected error, got nil")
	}
}

// Integration test

func TestIntegration_FullWorkflow(t *testing.T) {
	// Создаём конфигурацию через builder
	cfg := NewConfig(
		ConfigWithWidth(300),
		ConfigWithHeight(300),
		ConfigWithFontSize(45),
		ConfigWithOutputPath(t.TempDir()),
		ConfigWithOutputName("integration_test"),
		ConfigWithEmailTruncateLength(2),
	)

	// Создаём аватар
	a, err := New(cfg)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	result, err := a.Generate("integration.test@example.com")
	if err != nil {
		t.Fatalf("Generate() unexpected error: %v", err)
	}

	// Проверяем результат
	if result.FilePath == "" {
		t.Fatal("FilePath is empty")
	}
	if result.BackgroundColor == "" {
		t.Fatal("BackgroundColor is empty")
	}
	// integration — 11 символов, truncate=2, значит "in..."
	if result.Text != "in..." {
		t.Errorf("Text=%s, want 'in...'", result.Text)
	}

	// Проверяем что файл существует и не пустой
	stat, err := os.Stat(result.FilePath)
	if err != nil {
		t.Fatalf("file does not exist: %v", err)
	}
	if stat.Size() == 0 {
		t.Error("generated file is empty")
	}
}
