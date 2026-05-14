package config

import (
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.Width != 300 {
		t.Errorf("expected Width=300, got %d", cfg.Width)
	}
	if cfg.Height != 300 {
		t.Errorf("expected Height=300, got %d", cfg.Height)
	}
	if cfg.FontSize != 45 {
		t.Errorf("expected FontSize=45, got %f", cfg.FontSize)
	}
	if cfg.OutputPath != "." {
		t.Errorf("expected OutputPath='.', got %s", cfg.OutputPath)
	}
	if cfg.OutputName != "avatar" {
		t.Errorf("expected OutputName='avatar', got %s", cfg.OutputName)
	}
	if cfg.EmailTruncateLength != 2 {
		t.Errorf("expected EmailTruncateLength=2, got %d", cfg.EmailTruncateLength)
	}
	if cfg.FontPath != "" {
		t.Errorf("expected FontPath='', got %s", cfg.FontPath)
	}
	if cfg.ColorsPath != "" {
		t.Errorf("expected ColorsPath='', got %s", cfg.ColorsPath)
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := Default()
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestValidate_InvalidWidth(t *testing.T) {
	cfg := Default()
	cfg.Width = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for Width=0, got nil")
	}

	cfg.Width = -1
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for Width=-1, got nil")
	}
}

func TestValidate_InvalidHeight(t *testing.T) {
	cfg := Default()
	cfg.Height = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for Height=0, got nil")
	}
}

func TestValidate_InvalidFontSize(t *testing.T) {
	cfg := Default()
	cfg.FontSize = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for FontSize=0, got nil")
	}

	cfg.FontSize = -10
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for FontSize=-10, got nil")
	}
}

func TestValidate_InvalidEmailTruncateLength(t *testing.T) {
	cfg := Default()
	cfg.EmailTruncateLength = -1
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for EmailTruncateLength=-1, got nil")
	}
}

func TestValidate_CustomValues(t *testing.T) {
	cfg := &Config{
		Width:               500,
		Height:              500,
		FontSize:            60,
		FontPath:            "/path/to/font.ttf",
		OutputPath:          "/output",
		OutputName:          "custom",
		ColorsPath:          "/path/to/colors.json",
		EmailTruncateLength: 3,
	}

	if err := cfg.Validate(); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}

	if cfg.Width != 500 {
		t.Errorf("expected Width=500, got %d", cfg.Width)
	}
	if cfg.FontSize != 60 {
		t.Errorf("expected FontSize=60, got %f", cfg.FontSize)
	}
	if cfg.EmailTruncateLength != 3 {
		t.Errorf("expected EmailTruncateLength=3, got %d", cfg.EmailTruncateLength)
	}
}
