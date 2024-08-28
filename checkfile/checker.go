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
	FileAddress         string
	LogAddress          string
	Key                 []byte
	Directories         []string
	SupportedExtensions []string
	InfoJson            string
}

// CheckFile - метод для перевірки файлів у зазначених директоріях.
func (c *Checker) CheckFile() {

	// Завантажуємо список помилок
	errorPaths, err := logging.LoadErrorPaths()
	if err != nil {
		logging.Now().PrintLog(c.LogAddress, "[CheckFile] Помилка завантаження списку помилок", err.Error())
		return
	}

	// Проходження по всіх вказаних директоріях
	for _, dir := range c.Directories {
		if logging.IsPathInErrorList(dir, errorPaths) {
			//[CheckFile] Пропуск шляху з попередньою помилкою доступу"
			continue
		}

		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				logging.Now().PrintLog(c.LogAddress,
					"[CheckFile] Помилка доступу до шляху",
					"{Path :{"+path+"} Err :{"+err.Error()+"}}")

				logging.AddErrorPath(path, err.Error(), errorPaths)
				logging.SaveErrorPaths(errorPaths) // Зберігаємо оновлений список помилок

				return nil
			}

			// Перевіряє, чи підтримується тип файлу
			if !info.IsDir() && isSupportedFileType(path, c.SupportedExtensions) {
				changed, errorFileInfo := NewFileInfo().CheckAndWriteHash(path, "hashes.json")
				if errorFileInfo != nil {
					logging.Now().PrintLog(c.LogAddress,
						"[CheckFile] Помилка під час перевірки підтримування тип файлу", path)
				} else if changed {
					// Хеш файлу змінився
					sendfile.NewFILESender().SenderFile(c.FileAddress, c.LogAddress, path, c.Key, c.InfoJson)
				}
			}
			return nil
		})

		if err != nil {
			logging.Now().PrintLog(c.LogAddress,
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
