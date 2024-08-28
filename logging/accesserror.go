package logging

import (
	"encoding/json"
	"os"
)

type ErrorPath struct {
	Path  string `json:"path"`
	Error string `json:"error"`
}

type ErrorPaths struct {
	Paths []ErrorPath `json:"paths"`
}

const errorFilePath = "error_paths.json"

// Завантаження помилок з JSON-файлу
func LoadErrorPaths() (*ErrorPaths, error) {
	file, err := os.Open(errorFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &ErrorPaths{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var errorPaths ErrorPaths
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&errorPaths); err != nil {
		return nil, err
	}

	return &errorPaths, nil
}

// Збереження помилок до JSON-файлу
func SaveErrorPaths(errorPaths *ErrorPaths) error {
	file, err := os.Create(errorFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(errorPaths)
}

// Перевірка, чи містить список шляху помилок певний шлях
func IsPathInErrorList(path string, errorPaths *ErrorPaths) bool {
	for _, ep := range errorPaths.Paths {
		if ep.Path == path {
			return true
		}
	}
	return false
}

// Додавання нового шляху помилки до списку
func AddErrorPath(path, errorMsg string, errorPaths *ErrorPaths) {
	errorPaths.Paths = append(errorPaths.Paths, ErrorPath{
		Path:  path,
		Error: errorMsg,
	})
}
