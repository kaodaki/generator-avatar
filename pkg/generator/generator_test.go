package generator

import (
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/avatar-generator/avatar-generator/pkg/config"
)

func TestNew_WithDefaultConfig(t *testing.T) {
	cfg := config.Default()
	cfg.OutputPath = t.TempDir()

	gen, err := New(cfg)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}
	if gen == nil {
		t.Fatal("New() returned nil")
	}
}

func TestNew_WithNilConfig(t *testing.T) {
	gen, err := New(nil)
	if err != nil {
		t.Fatalf("New(nil) unexpected error: %v", err)
	}
	if gen == nil {
		t.Fatal("New(nil) returned nil")
	}
}

func TestNew_InvalidConfig(t *testing.T) {
	cfg := config.Default()
	cfg.Width = 0

	_, err := New(cfg)
	if err == nil {
		t.Error("New() with invalid config expected error, got nil")
	}
}

func TestGenerate(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := config.Default()
	cfg.OutputPath = tmpDir
	cfg.OutputName = "test_avatar"
	cfg.EmailTruncateLength = 100 // Большое значение чтобы не сокращать

	gen, err := New(cfg)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	result, err := gen.Generate("test@example.com")
	if err != nil {
		t.Fatalf("Generate() unexpected error: %v", err)
	}

	// Проверяем путь к файлу
	expectedPath := filepath.Join(tmpDir, "test_avatar.png")
	if result.FilePath != expectedPath {
		t.Errorf("FilePath=%s, want %s", result.FilePath, expectedPath)
	}

	// Проверяем что файл существует
	if _, err := os.Stat(result.FilePath); os.IsNotExist(err) {
		t.Errorf("file does not exist: %s", result.FilePath)
	}

	// Проверяем что цвет не пустой
	if result.BackgroundColor == "" {
		t.Error("BackgroundColor is empty")
	}

	// Проверяем текст (первая буква)
	if result.Text != "t" {
		t.Errorf("Text=%s, want 't'", result.Text)
	}

	// Проверяем что файл — валидный PNG
	data, err := os.ReadFile(result.FilePath)
	if err != nil {
		t.Fatalf("failed to read generated file: %v", err)
	}

	// PNG magic number: 89 50 4E 47 0D 0A 1A 0A
	if len(data) < 8 || data[0] != 0x89 || data[1] != 0x50 {
		t.Error("generated file is not a valid PNG")
	}

	// Проверяем что PNG можно декодировать
	file, err := os.Open(result.FilePath)
	if err != nil {
		t.Fatalf("failed to open generated file: %v", err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	// Проверяем размеры
	bounds := img.Bounds()
	if bounds.Dx() != 300 || bounds.Dy() != 300 {
		t.Errorf("image size=%dx%d, want 300x300", bounds.Dx(), bounds.Dy())
	}
}

func TestGenerate_WithTruncate(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := config.Default()
	cfg.OutputPath = tmpDir
	cfg.OutputName = "truncate_test"

	gen, err := New(cfg)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	// Email длиннее 2 символов — должен быть сокращён
	result, err := gen.Generate("verylongemail@example.com")
	if err != nil {
		t.Fatalf("Generate() unexpected error: %v", err)
	}

	if result.Text != "ve..." {
		t.Errorf("Text=%s, want 've...'", result.Text)
	}
}

func TestGenerate_CustomSize(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := config.Default()
	cfg.OutputPath = tmpDir
	cfg.OutputName = "custom_size"
	cfg.Width = 500
	cfg.Height = 500
	cfg.FontSize = 60

	gen, err := New(cfg)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	result, err := gen.Generate("user@example.com")
	if err != nil {
		t.Fatalf("Generate() unexpected error: %v", err)
	}

	// Проверяем размеры изображения
	file, _ := os.Open(result.FilePath)
	defer file.Close()
	img, _ := png.Decode(file)
	bounds := img.Bounds()

	if bounds.Dx() != 500 || bounds.Dy() != 500 {
		t.Errorf("image size=%dx%d, want 500x500", bounds.Dx(), bounds.Dy())
	}
}

func TestGenerate_EmptyEmail(t *testing.T) {
	cfg := config.Default()
	cfg.OutputPath = t.TempDir()

	gen, err := New(cfg)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	_, err = gen.Generate("")
	if err == nil {
		t.Error("Generate(\"\") expected error, got nil")
	}
}

func TestGenerate_ShortEmail(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := config.Default()
	cfg.OutputPath = tmpDir
	cfg.EmailTruncateLength = 5

	gen, err := New(cfg)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	// Email короче truncate length — показываем первую букву
	result, err := gen.Generate("ab@x")
	if err != nil {
		t.Fatalf("Generate() unexpected error: %v", err)
	}

	if result.Text != "a" {
		t.Errorf("Text=%s, want 'a'", result.Text)
	}
}

func TestGenerate_MultipleCalls(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := config.Default()
	cfg.OutputPath = tmpDir

	gen, err := New(cfg)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	// Генерируем несколько аватаров
	emails := []string{
		"user1@example.com",
		"user2@example.com",
		"admin@test.org",
		"support@company.io",
	}

	for _, email := range emails {
		cfg.OutputName = "avatar_" + strings.Split(email, "@")[0]
		gen.SetOutputName(cfg.OutputName)

		result, err := gen.Generate(email)
		if err != nil {
			t.Fatalf("Generate(%q) unexpected error: %v", email, err)
		}

		if _, err := os.Stat(result.FilePath); os.IsNotExist(err) {
			t.Errorf("file does not exist: %s", result.FilePath)
		}
	}
}

func TestSetOutputPath(t *testing.T) {
	cfg := config.Default()
	cfg.OutputPath = t.TempDir()

	gen, _ := New(cfg)

	newPath := filepath.Join(t.TempDir(), "subfolder")
	gen.SetOutputPath(newPath)
	gen.SetOutputName("path_test")

	result, err := gen.Generate("test@example.com")
	if err != nil {
		t.Fatalf("Generate() unexpected error: %v", err)
	}

	if !strings.HasPrefix(result.FilePath, newPath) {
		t.Errorf("FilePath=%s should start with %s", result.FilePath, newPath)
	}
}

func TestUpdateConfig(t *testing.T) {
	cfg := config.Default()
	cfg.OutputPath = t.TempDir()

	gen, _ := New(cfg)

	newCfg := config.Default()
	newCfg.OutputPath = t.TempDir()
	newCfg.Width = 400
	newCfg.Height = 400

	err := gen.UpdateConfig(newCfg)
	if err != nil {
		t.Fatalf("UpdateConfig() unexpected error: %v", err)
	}

	gen.SetOutputName("update_test")
	result, err := gen.Generate("test@example.com")
	if err != nil {
		t.Fatalf("Generate() unexpected error: %v", err)
	}

	// Проверяем размеры
	file, _ := os.Open(result.FilePath)
	defer file.Close()
	img, _ := png.Decode(file)
	bounds := img.Bounds()

	if bounds.Dx() != 400 || bounds.Dy() != 400 {
		t.Errorf("image size=%dx%d, want 400x400", bounds.Dx(), bounds.Dy())
	}
}

func TestUpdateConfig_Invalid(t *testing.T) {
	cfg := config.Default()
	cfg.OutputPath = t.TempDir()

	gen, _ := New(cfg)

	invalidCfg := config.Default()
	invalidCfg.Width = 0

	err := gen.UpdateConfig(invalidCfg)
	if err == nil {
		t.Error("UpdateConfig() with invalid config expected error, got nil")
	}
}

func TestPaletteCount(t *testing.T) {
	cfg := config.Default()
	cfg.OutputPath = t.TempDir()

	gen, _ := New(cfg)

	count := gen.PaletteCount()
	if count != 12 {
		t.Errorf("PaletteCount()=%d, want 12", count)
	}
}

func TestGenerate_WithCustomColorsFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Создаём кастомный файл цветов
	colorsFile := filepath.Join(tmpDir, "custom_colors.json")
	jsonData := `[
		{"name": "Custom Red", "hex": "#FF0000", "r": "255", "g": "0", "b": "0"},
		{"name": "Custom Green", "hex": "#00FF00", "r": "0", "g": "255", "b": "0"}
	]`
	os.WriteFile(colorsFile, []byte(jsonData), 0644)

	cfg := config.Default()
	cfg.OutputPath = tmpDir
	cfg.ColorsPath = colorsFile
	cfg.OutputName = "custom_colors_test"

	gen, err := New(cfg)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	if gen.PaletteCount() != 2 {
		t.Errorf("PaletteCount()=%d, want 2", gen.PaletteCount())
	}

	result, err := gen.Generate("test@example.com")
	if err != nil {
		t.Fatalf("Generate() unexpected error: %v", err)
	}

	// Проверяем что цвет один из наших
	if result.BackgroundColor != "#FF0000" && result.BackgroundColor != "#00FF00" {
		t.Errorf("Unexpected color: %s", result.BackgroundColor)
	}
}
