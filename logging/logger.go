package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// LogEntry представляє один запис логування
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

// Logger представляє логер для запису логів у JSON файл
type Logger struct {
	FilePath string
}

// NewLogger створює новий екземпляр Logger
func NewLogger(filePath string) *Logger {
	return &Logger{FilePath: filePath}
}

// Log записує новий запис у JSON файл
func (l *Logger) Log(level, message string) error {
	// Створення нового запису логування
	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level,
		Message:   message,
	}

	// Отримання існуючих записів
	var existingEntries []LogEntry
	if _, err := os.Stat(l.FilePath); !os.IsNotExist(err) {
		file, err := os.OpenFile(l.FilePath, os.O_RDWR, 0666)
		if err != nil {
			return err
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		err = decoder.Decode(&existingEntries)
		if err != nil && err != io.EOF {
			return err
		}
	}

	// Додавання нового запису
	existingEntries = append(existingEntries, entry)

	// Запис оновлених записів у файл
	file, err := os.Create(l.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Форматування JSON для кращої читабельності
	err = encoder.Encode(existingEntries)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	logger := NewLogger("logs.json")

	err := logger.Log("INFO", "This is an info message.")
	if err != nil {
		fmt.Printf("Error logging message: %v\n", err)
	}

	err = logger.Log("ERROR", "This is an error message.")
	if err != nil {
		fmt.Printf("Error logging message: %v\n", err)
	}
}
