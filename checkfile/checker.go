package checkfile

import (
	"Anthophila/sendfile"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Checker - структура, що містить дані для перевірки файлів.
type Checker struct {
	Address             string
	Key                 []byte
	Directories         []string
	SupportedExtensions []string
}

// NewChecker - конструктор для створення нового Checker.
func NewChecker(address string, key []byte, directories []string) *Checker {
	return &Checker{
		Address:     address,
		Key:         key,
		Directories: directories,
	}
}

// Checkfile - метод для перевірки файлів у зазначених директоріях.
func (c *Checker) Checkfile() error {
	// Проходження по всіх вказаних директоріях
	for _, dir := range c.Directories {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("Помилка доступу до шляху %s: %v\n", path, err)
				return nil
			}
			//Перевіряє, чи підтримується тип файлу
			if !info.IsDir() && isSupportedFileType(path, c.SupportedExtensions) {
				changed, errorFileInfo := NewFileInfo().CheckAndWriteHash(path, "hashes.json")
				if errorFileInfo != nil {
					fmt.Println("Помилка:", errorFileInfo)
				} else if changed {
					fmt.Println("Хеш файлу змінився")
					sender := sendfile.NewFILESender()
					senderError := sender.SenderFile(c.Address, path, c.Key)

					if senderError != nil {
						fmt.Println("Сталося помилка ", senderError)
					}
				} else {
					fmt.Println("Перевірка пошук та відправка нових і змінних файлі...")
				}

			}
			return nil
		})

		if err != nil {
			fmt.Printf("Помилка обходу шляху %s: %v\n", dir, err)
			return err
		}
	}
	return nil
}

/*
*	Функція isSupportedFileType:
*	Перевіряє, чи підтримується тип файлу. Повертає true, якщо розширення файлу є одним з підтримуваних
 */
func isSupportedFileType(file string, supportedExtensions []string) bool {
	// Перевіряємо, чи файл закінчується на ".enc"

	// Перевіряємо, чи файл має одне з підтримуваних розширень
	for _, ext := range supportedExtensions {
		if strings.HasSuffix(file, ext) {
			return true // Повертаємо true, якщо файл має підтримуване розширення
		}
	}
	return false // Повертаємо false, якщо файл не має підтримуваного розширення
}
