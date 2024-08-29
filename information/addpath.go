package information

import (
	"encoding/json"
	"io"
	"os"
)

type AddPath struct {
	// ... (поля для зберігання стану, якщо потрібні)
	FilePaths []string `json:"filePaths"`
}

func NewAddPath() *AddPath {
	return &AddPath{}
}
func (f *AddPath) AddFilePath(filePath string, jsonFilePath string) error {
	// Читання існуючого JSON-файлу
	file, err := os.OpenFile(jsonFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Декодування JSON
	decoder := json.NewDecoder(file)
	var store AddPath
	if err := decoder.Decode(&store); err != nil && err != io.EOF {
		return err
	}

	// Додавання нового посилання
	store.FilePaths = append(store.FilePaths, filePath)

	// Перезапис JSON-файлу
	file.Seek(0, 0) // Переміщення курсора на початок файлу
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(&store)
}
