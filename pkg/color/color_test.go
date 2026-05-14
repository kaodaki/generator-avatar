package color

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseHex_Valid(t *testing.T) {
	tests := []struct {
		hex      string
		expected string
		r        uint8
		g        uint8
		b        uint8
	}{
		{"#E74C3C", "#E74C3C", 231, 76, 60},
		{"#3498DB", "#3498DB", 52, 152, 219},
		{"#FFFFFF", "#FFFFFF", 255, 255, 255},
		{"#000000", "#000000", 0, 0, 0},
		{"E74C3C", "#E74C3C", 231, 76, 60},
		{"#e74c3c", "#E74C3C", 231, 76, 60},
	}

	for _, tt := range tests {
		c, err := ParseHex(tt.hex)
		if err != nil {
			t.Errorf("ParseHex(%q) unexpected error: %v", tt.hex, err)
			continue
		}
		if c.Hex != tt.expected {
			t.Errorf("ParseHex(%q) Hex=%s, want %s", tt.hex, c.Hex, tt.expected)
		}
		if c.R != tt.r || c.G != tt.g || c.B != tt.b {
			t.Errorf("ParseHex(%q) RGB=(%d,%d,%d), want (%d,%d,%d)",
				tt.hex, c.R, c.G, c.B, tt.r, tt.g, tt.b)
		}
	}
}

func TestParseHex_Invalid(t *testing.T) {
	invalidHexes := []string{
		"",
		"#",
		"#FFF",
		"#GGGGGG",
		"invalid",
		"#E74C3",
	}

	for _, h := range invalidHexes {
		_, err := ParseHex(h)
		if err == nil {
			t.Errorf("ParseHex(%q) expected error, got nil", h)
		}
	}
}

func TestContrastColor(t *testing.T) {
	tests := []struct {
		name        string
		hex         string
		expectedHex string
	}{
		{"White text on dark", "#34495E", "#FFFFFF"},
		{"White text on red", "#C0392B", "#FFFFFF"},
		{"Black text on light", "#FFFFFF", "#000000"},
		{"Black text on blue", "#3498DB", "#000000"},
		{"Black text on green", "#2ECC71", "#000000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := ParseHex(tt.hex)
			contrast := c.ContrastColor()
			if contrast.Hex != tt.expectedHex {
				t.Errorf("ContrastColor(%s)=%s, want %s", tt.hex, contrast.Hex, tt.expectedHex)
			}
		})
	}
}

func TestContrastColor_BrightnessThreshold(t *testing.T) {
	// Проверяем граничные значения яркости (YIQ = 128)
	// Цвет #808080 (серый) — примерно на границе
	c, _ := ParseHex("#808080")
	contrast := c.ContrastColor()
	// Для YIQ >= 128 должен быть чёрный, иначе белый
	yiq := ((int(c.R) * 299) + (int(c.G) * 587) + (int(c.B) * 114)) / 1000
	if yiq >= 128 && contrast.Hex != "#000000" {
		t.Errorf("Expected black for YIQ=%d, got %s", yiq, contrast.Hex)
	}
	if yiq < 128 && contrast.Hex != "#FFFFFF" {
		t.Errorf("Expected white for YIQ=%d, got %s", yiq, contrast.Hex)
	}
}

func TestLoadDefault(t *testing.T) {
	palette := LoadDefault()
	if palette == nil {
		t.Fatal("LoadDefault() returned nil")
	}

	count := palette.Count()
	if count != 12 {
		t.Errorf("expected 12 colors, got %d", count)
	}

	colors := palette.Colors()
	if len(colors) != 12 {
		t.Errorf("expected 12 colors from Colors(), got %d", len(colors))
	}
}

func TestPaletteRandom(t *testing.T) {
	palette := LoadDefault()

	// Проверяем что Random возвращает не-nil
	c := palette.Random()
	if c == nil {
		t.Fatal("Random() returned nil")
	}

	// Проверяем что HEX не пустой
	if c.Hex == "" {
		t.Error("Random() returned color with empty Hex")
	}

	// Проверяем что RGB компоненты валидны
	if c.R > 255 || c.G > 255 || c.B > 255 {
		t.Error("Random() returned color with invalid RGB values")
	}
}

func TestLoadFromFile(t *testing.T) {
	// Создаём временный JSON-файл
	tmpDir := t.TempDir()
	colorsFile := filepath.Join(tmpDir, "test_colors.json")

	jsonData := `[
		{"name": "Test Red", "hex": "#FF0000", "r": "255", "g": "0", "b": "0"},
		{"name": "Test Green", "hex": "#00FF00", "r": "0", "g": "255", "b": "0"},
		{"name": "Test Blue", "hex": "#0000FF", "r": "0", "g": "0", "b": "255"}
	]`

	if err := os.WriteFile(colorsFile, []byte(jsonData), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	palette, err := LoadFromFile(colorsFile)
	if err != nil {
		t.Fatalf("LoadFromFile(%q) unexpected error: %v", colorsFile, err)
	}

	if palette.Count() != 3 {
		t.Errorf("expected 3 colors, got %d", palette.Count())
	}
}

func TestLoadFromFile_NotFound(t *testing.T) {
	_, err := LoadFromFile("/nonexistent/path/colors.json")
	if err == nil {
		t.Error("LoadFromFile with nonexistent file expected error, got nil")
	}
}

func TestLoadFromFile_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	badFile := filepath.Join(tmpDir, "bad_colors.json")

	if err := os.WriteFile(badFile, []byte("not json"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	_, err := LoadFromFile(badFile)
	if err == nil {
		t.Error("LoadFromFile with invalid JSON expected error, got nil")
	}
}

func TestLoadFromFile_EmptyPalette(t *testing.T) {
	tmpDir := t.TempDir()
	emptyFile := filepath.Join(tmpDir, "empty_colors.json")

	if err := os.WriteFile(emptyFile, []byte("[]"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	_, err := LoadFromFile(emptyFile)
	if err == nil {
		t.Error("LoadFromFile with empty palette expected error, got nil")
	}
}

func TestToHex(t *testing.T) {
	c := &Color{Hex: "#E74C3C"}
	if c.ToHex() != "#E74C3C" {
		t.Errorf("ToHex()=%s, want #E74C3C", c.ToHex())
	}
}

func TestColors_Immutability(t *testing.T) {
	palette := LoadDefault()
	colors1 := palette.Colors()
	colors2 := palette.Colors()

	// Модифицируем первый результат
	colors1[0].Hex = "#MODIFIED"

	// Проверяем что второй не изменился
	if colors2[0].Hex == "#MODIFIED" {
		t.Error("Colors() returned mutable slice")
	}
}

func TestLuminance(t *testing.T) {
	// Белый — максимальная яркость
	white, _ := ParseHex("#FFFFFF")
	if white.Luminance() < 0.99 {
		t.Errorf("White luminance should be ~1.0, got %f", white.Luminance())
	}

	// Чёрный — минимальная яркость
	black, _ := ParseHex("#000000")
	if black.Luminance() > 0.01 {
		t.Errorf("Black luminance should be ~0.0, got %f", black.Luminance())
	}
}

func TestContrastRatio(t *testing.T) {
	white, _ := ParseHex("#FFFFFF")
	black, _ := ParseHex("#000000")

	// Контраст белого и чёрного должен быть 21:1
	ratio := white.ContrastRatio(black)
	if ratio < 20.5 || ratio > 21.5 {
		t.Errorf("White/Black contrast ratio expected ~21, got %f", ratio)
	}
}
