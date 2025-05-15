package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

var (
	Dir  string
	Port int
	help bool
)

var (
	ErrorLogPath = "logs/error.log"
	InfoLogPath  = "logs/info.log"
	DebugLogPath = "logs/debug.log"
)

var helpTxt = `
Coffee Shop Management System

Usage:
  hot-coffee [--port <N>] [--dir <S>] 
  hot-coffee --help

Options:
  --help       Show this screen.
  --port N     Port number.
  --dir S      Path to the data directory
`

var usageTxt = `
Usage:
  hot-coffee [--port <N>] [--dir <S>] 
  hot-coffee --help

Options:
  --help       Show help.
  --port N     Port number.
  --dir S      Path to the data directory
`

// Инициализация флагов командной строки
func init() {
	flag.Usage = func() {
		fmt.Printf(usageTxt)
	}

	flag.StringVar(&Dir, "dir", "./data", "Path to the data directory.")
	flag.IntVar(&Port, "port", 8080, "Port number.")
	flag.BoolVar(&help, "help", false, "Show help.")
	flag.Parse()

	// Если задан флаг --help, выводим справку и выходим
	if help || slices.Contains(os.Args, "--help") {
		fmt.Print(helpTxt)
		os.Exit(0)
	}
}

// Инициализация конфигурации
func InitConfig() error {
	// Проверяем корректность флагов
	if err := checkFlags(); err != nil {
		return err
	}

	// Создаем файлы, если их еще нет
	if err := createFileIfNotExists("menu.json"); err != nil {
		return err
	}
	if err := createFileIfNotExists("order.json"); err != nil {
		return err
	}
	if err := createFileIfNotExists("inventory.json"); err != nil {
		return err
	}

	return nil
}

// Проверяем корректность флагов
func checkFlags() error {
	// Проверяем наличие неожиданных аргументов
	if len(flag.Args()) > 0 {
		return fmt.Errorf("Unexpected arguments: %v\n%s", flag.Args(), usageTxt)
	}

	// Проверяем корректность номера порта
	if Port < 1024 || Port > 49151 {
		return fmt.Errorf("Port number must be in the range [1024, 49151]\n%s", usageTxt)
	}

	// Проверяем существование и доступность директории с данными
	if stat, err := os.Stat(Dir); err != nil || !stat.IsDir() {
		return fmt.Errorf("Invalid data directory: %s\n%s", Dir, usageTxt)
	}

	return nil
}

// Создаем файл, если он еще не существует
func createFileIfNotExists(fileName string) error {
	path := filepath.Join(Dir, fileName)
	if !hasFile(path) {
		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create %s: %v", fileName, err)
		}

		_, err = file.Write([]byte("[]"))
		if err != nil {
			return err
		}

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				fmt.Printf("failed to close file: %v", err)
			}
		}(file)
	}
	return nil
}

// Проверяем, существует ли файл
func hasFile(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
