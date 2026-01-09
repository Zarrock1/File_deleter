package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	fmt.Println("=== УДАЛЕНИЕ ФАЙЛОВ С (ЦИФРА) В НАЗВАНИИ ===")
	fmt.Println()
	
	// Если путь передан как аргумент
	var folderPath string
	if len(os.Args) > 1 {
		folderPath = os.Args[1]
		fmt.Printf("Путь из аргумента: %s\n", folderPath)
	} else {
		// Запрашиваем путь
		fmt.Print("Введите путь к папке: ")
		fmt.Scanln(&folderPath)
	}

	// Убираем лишние пробелы
	folderPath = strings.TrimSpace(folderPath)
	
	if folderPath == "" {
		fmt.Println("Ошибка: путь не указан")
		fmt.Println("Использование:")
		fmt.Println("  1. Перетащите папку на программу")
		fmt.Println("  2. Или запустите: program.exe \"C:\\Папка\\с пробелами\"")
		return
	}

	// Если путь в кавычках, убираем их
	if strings.HasPrefix(folderPath, "\"") && strings.HasSuffix(folderPath, "\"") {
		folderPath = folderPath[1 : len(folderPath)-1]
	}

	// Получаем абсолютный путь
	absPath, err := filepath.Abs(folderPath)
	if err != nil {
		fmt.Printf("Ошибка обработки пути: %v\n", err)
		return
	}

	// Проверяем существование папки
	info, err := os.Stat(absPath)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		fmt.Printf("Путь: %s\n", absPath)
		return
	}
	
	if !info.IsDir() {
		fmt.Printf("Ошибка: %s не является папкой\n", absPath)
		return
	}

	fmt.Printf("\n✅ Найдена папка: %s\n", absPath)
	
	// Подтверждение
	fmt.Print("\n⚠️  УДАЛИТЬ все файлы с (1), (2) и т.д.? (y/N): ")
	
	var confirm string
	fmt.Scanln(&confirm)
	confirm = strings.TrimSpace(strings.ToLower(confirm))
	
	if confirm != "y" && confirm != "yes" && confirm != "д" && confirm != "да" {
		fmt.Println("❌ Отменено")
		return
	}

	fmt.Println("\n🔍 Начинаю сканирование и удаление...")
	fmt.Println("══════════════════════════════════════════")

	deletedCount := deleteFiles(absPath)

	fmt.Println("══════════════════════════════════════════")
	fmt.Printf("✅ Готово! Удалено файлов: %d\n", deletedCount)
	
	// Ждем нажатия Enter перед выходом
	fmt.Print("\nНажмите Enter для выхода...")
	fmt.Scanln()
}

func deleteFiles(folderPath string) int {
	// Регулярное выражение для поиска (цифра)
	pattern := regexp.MustCompile(`\(\d+\)`)
	deletedCount := 0
	
	// Рекурсивно обходим папку
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Пропускаем файлы с ошибками доступа
			return nil
		}
		
		// Пропускаем папки
		if info.IsDir() {
			return nil
		}
		
		// Проверяем имя файла
		filename := info.Name()
		if pattern.MatchString(filename) {
			// Пытаемся удалить файл
			err := os.Remove(path)
			if err != nil {
				fmt.Printf("❌ Ошибка: %s - %v\n", filename, err)
			} else {
				fmt.Printf("✅ Удален: %s\n", filename)
				deletedCount++
			}
		}
		
		return nil
	})
	
	if err != nil {
		fmt.Printf("⚠️  Внимание при сканировании: %v\n", err)
	}
	
	return deletedCount
}