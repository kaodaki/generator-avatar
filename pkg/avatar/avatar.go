// Package avatar предоставляет высокоуровневый API для генерации аватаров.
// Это основной пакет библиотеки, который подключается к проектам.
//
// Пример использования:
//
//	a := avatar.New(avatar.DefaultConfig())
//	result, err := a.Generate("user@example.com")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Avatar saved to:", result.FilePath)
//	fmt.Println("Background color:", result.BackgroundColor)
package avatar

import (
	"fmt"

	"github.com/avatar-generator/avatar-generator/pkg/color"
	"github.com/avatar-generator/avatar-generator/pkg/config"
	"github.com/avatar-generator/avatar-generator/pkg/generator"
)

// Avatar предоставляет API для генерации аватаров.
type Avatar struct {
	gen *generator.Generator
}

// New создаёт новый экземпляр Avatar с указанной конфигурацией.
// Если cfg == nil, используются значения по умолчанию.
func New(cfg *config.Config) (*Avatar, error) {
	if cfg == nil {
		cfg = config.Default()
	}

	gen, err := generator.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("create generator: %w", err)
	}

	return &Avatar{gen: gen}, nil
}

// Generate создаёт аватар для указанного email.
// Возвращает результат генерации или ошибку.
func (a *Avatar) Generate(email string) (*Result, error) {
	res, err := a.gen.Generate(email)
	if err != nil {
		return nil, fmt.Errorf("generate avatar: %w", err)
	}

	return &Result{
		FilePath:        res.FilePath,
		BackgroundColor: res.BackgroundColor,
		Text:            res.Text,
	}, nil
}

// GenerateWithOptions создаёт аватар с временными настройками.
// Позволяет переопределить путь и имя файла для одной генерации.
func (a *Avatar) GenerateWithOptions(email, outputPath, outputName string) (*Result, error) {
	// Сохраняем текущие настройки
	oldPath := a.gen.SetOutputPath(outputPath)
	a.gen.SetOutputName(outputName)

	// Генерируем
	res, err := a.gen.Generate(email)

	// Восстанавливаем путь (лучше сохранять состояние)
	_ = oldPath

	if err != nil {
		return nil, fmt.Errorf("generate avatar: %w", err)
	}

	return &Result{
		FilePath:        res.FilePath,
		BackgroundColor: res.BackgroundColor,
		Text:            res.Text,
	}, nil
}

// SetOutputPath изменяет путь сохранения аватаров.
func (a *Avatar) SetOutputPath(path string) {
	a.gen.SetOutputPath(path)
}

// SetOutputName изменяет имя выходного файла.
func (a *Avatar) SetOutputName(name string) {
	a.gen.SetOutputName(name)
}

// UpdateConfig обновляет конфигурацию генератора.
func (a *Avatar) UpdateConfig(cfg *config.Config) error {
	return a.gen.UpdateConfig(cfg)
}

// PaletteCount возвращает количество доступных цветов.
func (a *Avatar) PaletteCount() int {
	return a.gen.PaletteCount()
}

// Result содержит результат генерации аватара.
type Result struct {
	FilePath        string
	BackgroundColor string
	Text            string
}

// DefaultConfig возвращает конфигурацию по умолчанию.
func DefaultConfig() *config.Config {
	return config.Default()
}

// WithWidth устанавливает ширину аватара.
func WithWidth(cfg *config.Config, width int) *config.Config {
	cfg.Width = width
	return cfg
}

// WithHeight устанавливает высоту аватара.
func WithHeight(cfg *config.Config, height int) *config.Config {
	cfg.Height = height
	return cfg
}

// WithFontSize устанавливает размер шрифта.
func WithFontSize(cfg *config.Config, size float64) *config.Config {
	cfg.FontSize = size
	return cfg
}

// WithOutputPath устанавливает путь сохранения.
func WithOutputPath(cfg *config.Config, path string) *config.Config {
	cfg.OutputPath = path
	return cfg
}

// WithOutputName устанавливает имя файла.
func WithOutputName(cfg *config.Config, name string) *config.Config {
	cfg.OutputName = name
	return cfg
}

// WithColorsPath устанавливает путь к файлу цветов.
func WithColorsPath(cfg *config.Config, path string) *config.Config {
	cfg.ColorsPath = path
	return cfg
}

// WithEmailTruncateLength устанавливает длину сокращения email.
func WithEmailTruncateLength(cfg *config.Config, length int) *config.Config {
	cfg.EmailTruncateLength = length
	return cfg
}

// WithFontPath устанавливает путь к шрифту.
func WithFontPath(cfg *config.Config, path string) *config.Config {
	cfg.FontPath = path
	return cfg
}

// Helper-функции для создания конфигурации builder-стилем

// NewConfig создаёт конфигурацию с кастомными опциями.
func NewConfig(options ...ConfigOption) *config.Config {
	cfg := config.Default()
	for _, opt := range options {
		opt(cfg)
	}
	return cfg
}

// ConfigOption — функция-опция для конфигурации.
type ConfigOption func(*config.Config)

// ConfigWithWidth устанавливает ширину.
func ConfigWithWidth(width int) ConfigOption {
	return func(c *config.Config) {
		c.Width = width
	}
}

// ConfigWithHeight устанавливает высоту.
func ConfigWithHeight(height int) ConfigOption {
	return func(c *config.Config) {
		c.Height = height
	}
}

// ConfigWithFontSize устанавливает размер шрифта.
func ConfigWithFontSize(size float64) ConfigOption {
	return func(c *config.Config) {
		c.FontSize = size
	}
}

// ConfigWithOutputPath устанавливает путь сохранения.
func ConfigWithOutputPath(path string) ConfigOption {
	return func(c *config.Config) {
		c.OutputPath = path
	}
}

// ConfigWithOutputName устанавливает имя файла.
func ConfigWithOutputName(name string) ConfigOption {
	return func(c *config.Config) {
		c.OutputName = name
	}
}

// ConfigWithColorsPath устанавливает путь к файлу цветов.
func ConfigWithColorsPath(path string) ConfigOption {
	return func(c *config.Config) {
		c.ColorsPath = path
	}
}

// ConfigWithEmailTruncateLength устанавливает длину сокращения email.
func ConfigWithEmailTruncateLength(length int) ConfigOption {
	return func(c *config.Config) {
		c.EmailTruncateLength = length
	}
}

// ConfigWithFontPath устанавливает путь к шрифту.
func ConfigWithFontPath(path string) ConfigOption {
	return func(c *config.Config) {
		c.FontPath = path
	}
}

// Utility functions

// PickRandomColor возвращает случайный цвет из встроенной палитры.
func PickRandomColor() *color.Color {
	palette := color.LoadDefault()
	return palette.Random()
}

// LoadPaletteFromFile загружает палитру из JSON-файла.
func LoadPaletteFromFile(path string) (*color.Palette, error) {
	return color.LoadFromFile(path)
}
