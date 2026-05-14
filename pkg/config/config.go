// Package config предоставляет конфигурационные структуры для генератора аватаров.
// Все параметры могут быть изменены при подключении библиотеки к проекту.
package config

import "errors"

// Config содержит все настраиваемые параметры для генерации аватаров.
type Config struct {
	// Width — ширина аватара в пикселях (по умолчанию: 300)
	Width int

	// Height — высота аватара в пикселях (по умолчанию: 300)
	Height int

	// FontSize — размер шрифта для отображения текста (по умолчанию: 45)
	FontSize float64

	// FontPath — путь к файлу шрифта TTF
	// По умолчанию: встроенный шрифт Share Tech Mono
	FontPath string

	// OutputPath — директория для сохранения сгенерированных аватаров
	// По умолчанию: текущая директория
	OutputPath string

	// OutputName — имя выходного файла (без расширения)
	// По умолчанию: avatar
	OutputName string

	// ColorsPath — путь к JSON-файлу с палитрой цветов
	// По умолчанию: встроенный набор цветов
	ColorsPath string

	// EmailTruncateLength — количество символов email перед многоточием
	// По умолчанию: 2
	EmailTruncateLength int
}

// Validate проверяет корректность конфигурации.
func (c *Config) Validate() error {
	if c.Width <= 0 {
		return errors.New("width must be greater than 0")
	}
	if c.Height <= 0 {
		return errors.New("height must be greater than 0")
	}
	if c.FontSize <= 0 {
		return errors.New("font size must be greater than 0")
	}
	if c.EmailTruncateLength < 0 {
		return errors.New("email truncate length must be non-negative")
	}
	return nil
}

// Default возвращает конфигурацию со значениями по умолчанию.
func Default() *Config {
	return &Config{
		Width:               300,
		Height:              300,
		FontSize:            45,
		FontPath:            "",
		OutputPath:          ".",
		OutputName:          "avatar",
		ColorsPath:          "",
		EmailTruncateLength: 2,
	}
}
