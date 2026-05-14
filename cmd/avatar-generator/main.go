// avatar-generator — CLI-приложение для демонстрации генерации аватаров.
//
// Использование:
//
//	go run cmd/avatar-generator/main.go -email user@example.com
//	go run cmd/avatar-generator/main.go -email user@example.com -output ./avatars -name myavatar
//	go run cmd/avatar-generator/main.go -email user@example.com -width 500 -height 500 -fontsize 60
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/avatar-generator/avatar-generator/pkg/avatar"
	"github.com/avatar-generator/avatar-generator/pkg/config"
)

func main() {
	// Парсим флаги
	email := flag.String("email", "", "Email пользователя (обязательный)")
	outputPath := flag.String("output", "./output", "Путь для сохранения аватара")
	outputName := flag.String("name", "avatar", "Имя выходного файла (без .png)")
	width := flag.Int("width", 300, "Ширина аватара в пикселях")
	height := flag.Int("height", 300, "Высота аватара в пикселях")
	fontSize := flag.Float64("fontsize", 45, "Размер шрифта")
	fontPath := flag.String("font", "", "Путь к файлу шрифта TTF (опционально)")
	colorsPath := flag.String("colors", "", "Путь к JSON-файлу с цветами (опционально)")
	truncateLen := flag.Int("truncate", 2, "Длина сокращения email с многоточием")

	flag.Parse()

	if *email == "" {
		fmt.Fprintf(os.Stderr, "Ошибка: параметр -email обязательный\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Создаём конфигурацию
	cfg := config.Default()
	cfg.Width = *width
	cfg.Height = *height
	cfg.FontSize = *fontSize
	cfg.FontPath = *fontPath
	cfg.OutputPath = *outputPath
	cfg.OutputName = *outputName
	cfg.ColorsPath = *colorsPath
	cfg.EmailTruncateLength = *truncateLen

	// Создаём генератор
	gen, err := avatar.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка создания генератора: %v\n", err)
		os.Exit(1)
	}

	// Генерируем аватар
	result, err := gen.Generate(*email)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка генерации аватара: %v\n", err)
		os.Exit(1)
	}

	// Выводим результат
	fmt.Println("=" + repeatString("=", 50))
	fmt.Println("  Аватар успешно сгенерирован!")
	fmt.Println("=" + repeatString("=", 50))
	fmt.Printf("  Файл:       %s\n", result.FilePath)
	fmt.Printf("  Цвет фона:  %s\n", result.BackgroundColor)
	fmt.Printf("  Текст:      %s\n", result.Text)
	fmt.Printf("  Размер:     %dx%d px\n", *width, *height)
	fmt.Printf("  Email:      %s\n", *email)
	fmt.Println("=" + repeatString("=", 50))
}

func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
