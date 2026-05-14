// Package color предоставляет типы и функции для работы с цветовой палитрой аватаров.
package color

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Color представляет цвет в различных форматах.
type Color struct {
	Name string `json:"name"`
	Hex  string `json:"hex"`
	R    uint8  `json:"-"`
	G    uint8  `json:"-"`
	B    uint8  `json:"-"`
}

// UnmarshalJSON реализует кастомную десериализацию из JSON.
func (c *Color) UnmarshalJSON(data []byte) error {
	type rawColor struct {
		Name string `json:"name"`
		Hex  string `json:"hex"`
		R    string `json:"r"`
		G    string `json:"g"`
		B    string `json:"b"`
	}

	var raw rawColor
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("unmarshal color: %w", err)
	}

	c.Name = raw.Name
	c.Hex = raw.Hex

	// Парсим RGB компоненты
	r, err := strconv.Atoi(raw.R)
	if err != nil {
		return fmt.Errorf("parse R component %q: %w", raw.R, err)
	}
	g, err := strconv.Atoi(raw.G)
	if err != nil {
		return fmt.Errorf("parse G component %q: %w", raw.G, err)
	}
	b, err := strconv.Atoi(raw.B)
	if err != nil {
		return fmt.Errorf("parse B component %q: %w", raw.B, err)
	}

	c.R = uint8(r)
	c.G = uint8(g)
	c.B = uint8(b)

	return nil
}

// ToHex возвращает HEX-код цвета.
func (c *Color) ToHex() string {
	return c.Hex
}

// ContrastColor определяет контрастный цвет (белый или чёрный)
// на основе яркости фона по алгоритму YIQ.
func (c *Color) ContrastColor() *Color {
	yiq := ((int(c.R) * 299) + (int(c.G) * 587) + (int(c.B) * 114)) / 1000
	if yiq >= 128 {
		return &Color{Name: "Black", Hex: "#000000", R: 0, G: 0, B: 0}
	}
	return &Color{Name: "White", Hex: "#FFFFFF", R: 255, G: 255, B: 255}
}

// Palette представляет набор цветов для аватаров.
type Palette struct {
	colors []Color
}

// LoadFromFile загружает палитру цветов из JSON-файла.
func LoadFromFile(path string) (*Palette, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read colors file: %w", err)
	}

	var colors []Color
	if err := json.Unmarshal(data, &colors); err != nil {
		return nil, fmt.Errorf("parse colors JSON: %w", err)
	}

	if len(colors) == 0 {
		return nil, fmt.Errorf("color palette is empty")
	}

	return &Palette{colors: colors}, nil
}

// LoadDefault загружает встроенную палитру цветов.
func LoadDefault() *Palette {
	// Встроенный набор цветов (копия из colors.json)
	defaultColors := []Color{
		{Name: "Ализариновый красный", Hex: "#E74C3C", R: 231, G: 76, B: 60},
		{Name: "Питерский синий", Hex: "#3498DB", R: 52, G: 152, B: 219},
		{Name: "Изумрудный", Hex: "#2ECC71", R: 46, G: 204, B: 113},
		{Name: "Оранжевый", Hex: "#F39C12", R: 243, G: 156, B: 18},
		{Name: "Аметистовый", Hex: "#9B59B6", R: 155, G: 89, B: 182},
		{Name: "Бирюзовый", Hex: "#1ABC9C", R: 26, G: 188, B: 156},
		{Name: "Морковный", Hex: "#E67E22", R: 230, G: 126, B: 34},
		{Name: "Мокрый асфальт", Hex: "#34495E", R: 52, G: 73, B: 94},
		{Name: "Зелёное море", Hex: "#16A085", R: 22, G: 160, B: 133},
		{Name: "Тыквенный", Hex: "#D35400", R: 211, G: 84, B: 0},
		{Name: "Тёмно-фиолетовый", Hex: "#8E44AD", R: 142, G: 68, B: 173},
		{Name: "Тёмно-красный", Hex: "#C0392B", R: 192, G: 57, B: 43},
	}
	return &Palette{colors: defaultColors}
}

// Random выбирает случайный цвет из палитры.
func (p *Palette) Random() *Color {
	idx := rand.Intn(len(p.colors))
	chosen := p.colors[idx]
	return &Color{
		Name: chosen.Name,
		Hex:  chosen.Hex,
		R:    chosen.R,
		G:    chosen.G,
		B:    chosen.B,
	}
}

// Count возвращает количество цветов в палитре.
func (p *Palette) Count() int {
	return len(p.colors)
}

// Colors возвращает копию всех цветов палитры.
func (p *Palette) Colors() []Color {
	result := make([]Color, len(p.colors))
	copy(result, p.colors)
	return result
}

// ParseHex парсит HEX-строку в структуру Color.
func ParseHex(hex string) (*Color, error) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return nil, fmt.Errorf("invalid hex format: %s", hex)
	}

	r, err := strconv.ParseInt(hex[0:2], 16, 32)
	if err != nil {
		return nil, fmt.Errorf("parse red component: %w", err)
	}
	g, err := strconv.ParseInt(hex[2:4], 16, 32)
	if err != nil {
		return nil, fmt.Errorf("parse green component: %w", err)
	}
	b, err := strconv.ParseInt(hex[4:6], 16, 32)
	if err != nil {
		return nil, fmt.Errorf("parse blue component: %w", err)
	}

	return &Color{
		Hex: fmt.Sprintf("#%s", strings.ToUpper(hex)),
		R:   uint8(r),
		G:   uint8(g),
		B:   uint8(b),
	}, nil
}

// Luminance вычисляет относительную яркость цвета (0-1).
func (c *Color) Luminance() float64 {
	r := float64(c.R) / 255.0
	g := float64(c.G) / 255.0
	b := float64(c.B) / 255.0

	// Применяем гамма-коррекцию
	r = applyGammaCorrection(r)
	g = applyGammaCorrection(g)
	b = applyGammaCorrection(b)

	return 0.2126*r + 0.7152*g + 0.0722*b
}

// ContrastRatio вычисляет коэффициент контраста между двумя цветами.
func (c *Color) ContrastRatio(other *Color) float64 {
	l1 := c.Luminance()
	l2 := other.Luminance()

	if l1 < l2 {
		l1, l2 = l2, l1
	}

	return (l1 + 0.05) / (l2 + 0.05)
}

func applyGammaCorrection(c float64) float64 {
	if c <= 0.03928 {
		return c / 12.92
	}
	return math.Pow((c+0.055)/1.055, 2.4)
}
