// Package generator предоставляет функционал для генерации PNG-аватаров.
package generator

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	avatarColor "github.com/avatar-generator/avatar-generator/pkg/color"
	"github.com/avatar-generator/avatar-generator/pkg/config"
)

// Result содержит результат генерации аватара.
type Result struct {
	// FilePath — полный путь к сохранённому файлу
	FilePath string

	// BackgroundColor — HEX-код выбранного цвета фона
	BackgroundColor string

	// Text — текст, который был отображён на аватаре
	Text string
}

// Generator создаёт аватары с настраиваемыми параметрами.
type Generator struct {
	cfg     *config.Config
	palette *avatarColor.Palette
}

// New создаёт новый генератор с указанной конфигурацией.
func New(cfg *config.Config) (*Generator, error) {
	if cfg == nil {
		cfg = config.Default()
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	g := &Generator{
		cfg: cfg,
	}

	// Загружаем палитру цветов
	if err := g.loadPalette(); err != nil {
		return nil, fmt.Errorf("load palette: %w", err)
	}

	return g, nil
}

// loadPalette загружает палитру цветов из файла или использует встроенную.
func (g *Generator) loadPalette() error {
	if g.cfg.ColorsPath != "" {
		palette, err := avatarColor.LoadFromFile(g.cfg.ColorsPath)
		if err != nil {
			return err
		}
		g.palette = palette
		return nil
	}
	g.palette = avatarColor.LoadDefault()
	return nil
}

// Generate создаёт аватар на основе email.
// Возвращает путь к файлу и HEX-код выбранного цвета.
func (g *Generator) Generate(email string) (*Result, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	// Выбираем случайный цвет
	bgColor := g.palette.Random()
	textColor := bgColor.ContrastColor()

	// Определяем текст для отображения
	displayText := g.extractText(email)

	// Создаём изображение
	img := image.NewRGBA(image.Rect(0, 0, g.cfg.Width, g.cfg.Height))

	// Заливаем фон
	bgRGBA := toImageColor(bgColor)
	draw.Draw(img, img.Bounds(), &image.Uniform{C: bgRGBA}, image.Point{}, draw.Src)

	// Рисуем текст
	if err := g.drawText(img, displayText, textColor); err != nil {
		return nil, fmt.Errorf("draw text: %w", err)
	}

	// Формируем путь к файлу
	filePath := g.buildFilePath()

	// Создаём директорию если нужно
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create output directory: %w", err)
	}

	// Сохраняем PNG
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return nil, fmt.Errorf("encode png: %w", err)
	}

	return &Result{
		FilePath:        filePath,
		BackgroundColor: bgColor.ToHex(),
		Text:            displayText,
	}, nil
}

// extractText извлекает первую букву email или сокращает email с многоточием.
func (g *Generator) extractText(email string) string {
	email = strings.TrimSpace(email)

	// Получаем первую руну
	firstRune, size := utf8.DecodeRuneInString(email)
	if firstRune == utf8.RuneError && size == 0 {
		return "?"
	}

	_ = string(firstRune)

	// Если email длиннее допустимого — сокращаем
	emailLen := utf8.RuneCountInString(email)
	if g.cfg.EmailTruncateLength > 0 && emailLen > g.cfg.EmailTruncateLength {
		runes := []rune(email)
		truncateLen := g.cfg.EmailTruncateLength
		if truncateLen > len(runes) {
			truncateLen = len(runes)
		}
		return string(runes[:truncateLen]) + "..."
	}

	return string(firstRune)
}

// drawText рисует текст по центру изображения.
func (g *Generator) drawText(img *image.RGBA, text string, textColor *avatarColor.Color) error {
	// Вычисляем масштабный коэффициент на основе FontSize
	// Базовый шрифт 7x5 пикселей, масштабируем до нужного размера
	textRunes := []rune(text)

	// Вычисляем размеры текста в пикселях (без масштабирования)
	totalWidth := len(textRunes) * charWidth()
	// Добавляем пробелы между символами
	if len(textRunes) > 1 {
		totalWidth += (len(textRunes) - 1)
	}
	totalHeight := charHeight()

	// Вычисляем масштаб чтобы текст занимал ~60% от размера аватара
	desiredTextSize := float64(g.cfg.Width) * 0.6
	scale := 1
	if totalWidth > 0 {
		scale = int(desiredTextSize / float64(totalWidth))
	}
	if scale < 1 {
		scale = 1
	}

	// Вычисляем итоговые размеры
	scaledWidth := totalWidth * scale
	if len(textRunes) > 1 {
		scaledWidth = totalWidth*scale + (len(textRunes)-1)*scale
	}
	scaledHeight := totalHeight * scale

	// Позиция для центрирования
	startX := (g.cfg.Width - scaledWidth) / 2
	startY := (g.cfg.Height - scaledHeight) / 2

	txtClr := toImageColor(textColor)

	// Рисуем каждый символ
	for i, r := range textRunes {
		charData := getChar(r)
		charX := startX + i*(charWidth()*scale+scale)

		for row := 0; row < charHeight(); row++ {
			for col := 0; col < charWidth(); col++ {
				if charData[row][col] == 1 {
					// Рисуем масштабированный пиксель
					for sy := 0; sy < scale; sy++ {
						for sx := 0; sx < scale; sx++ {
							px := charX + col*scale + sx
							py := startY + row*scale + sy
							if px >= 0 && px < g.cfg.Width && py >= 0 && py < g.cfg.Height {
								img.Set(px, py, txtClr)
							}
						}
					}
				}
			}
		}
	}

	return nil
}

// buildFilePath формирует полный путь к выходному файлу.
func (g *Generator) buildFilePath() string {
	name := g.cfg.OutputName
	if name == "" {
		name = "avatar"
	}
	if !strings.HasSuffix(name, ".png") {
		name += ".png"
	}
	return filepath.Join(g.cfg.OutputPath, name)
}

// toImageColor конвертирует avatarColor.Color в color.Color.
func toImageColor(c *avatarColor.Color) color.Color {
	return color.RGBA{R: c.R, G: c.G, B: c.B, A: 255}
}

// UpdateConfig позволяет обновить конфигурацию генератора.
func (g *Generator) UpdateConfig(cfg *config.Config) error {
	if err := cfg.Validate(); err != nil {
		return err
	}
	g.cfg = cfg

	// Перезагружаем палитру если изменился путь к цветам
	if cfg.ColorsPath != "" {
		return g.loadPalette()
	}
	return nil
}

// SetOutputPath изменяет путь сохранения.
func (g *Generator) SetOutputPath(path string) string {
	old := g.cfg.OutputPath
	g.cfg.OutputPath = path
	return old
}

// SetOutputName изменяет имя выходного файла.
func (g *Generator) SetOutputName(name string) {
	g.cfg.OutputName = name
}

// PaletteCount возвращает количество цветов в палитре.
func (g *Generator) PaletteCount() int {
	return g.palette.Count()
}
