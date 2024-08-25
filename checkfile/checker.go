package checkfile

import (
	"Anthophila/logging"
	"Anthophila/sendfile"
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
	InfoJson            string
}

// CheckFile - метод для перевірки файлів у зазначених директоріях.
func (c *Checker) CheckFile() {

	// Проходження по всіх вказаних директоріях
	for _, dir := range c.Directories {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				logging.Now().PrintLog(
					"[CheckFile] Помилка доступу до шляху",
					"{Path :{"+path+"} Err :{"+err.Error()+"}}")
				return nil
			}
			//Перевіряє, чи підтримується тип файлу
			if !info.IsDir() && isSupportedFileType(path, c.SupportedExtensions) {
				changed, errorFileInfo := NewFileInfo().CheckAndWriteHash(path, "hashes.json")
				if errorFileInfo != nil {
					logging.Now().PrintLog(
						"[CheckFile] Помилка під час перевірки підтримування тип файлу", path)
				} else if changed {
					//Хеш файлу змінився
					sendfile.NewFILESender().SenderFile(c.Address, path, c.Key, c.InfoJson)
				} else {
					//fmt.Println("Перевірка пошук та відряджання нових і змінних файлі...")
				}

			}
			return nil
		})

		if err != nil {
			logging.Now().PrintLog(
				"[CheckFile] Помилка обходу шляху",
				"{Dir :{"+dir+"} Err :{"+err.Error()+"}}")
		}
	}
}

/*
*	Функція isSupportedFileType:
*	Перевіряє, чи підтримується тип файлу. Повертає true, якщо розширення файлу є одним з підтримуваних
 */
func isSupportedFileType(file string, supportedExtensions []string) bool {
	// Перевіряємо, чи файл має одне з підтримуваних розширень
	for _, ext := range supportedExtensions {
		if strings.HasSuffix(file, ext) {
			return true // Повертаємо true, якщо файл має підтримуване розширення
		}
	}
	return false // Повертаємо false, якщо файл не має підтримуваного розширення
}
