package checkfile

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FileInfo struct {
	Path string `json:"path"`
	Name string `json:"name"`
	Hash string `json:"hash"`
}

func NewFileInfo() *FileInfo {
	return &FileInfo{}
}
func (fi FileInfo) CheckAndWriteHash(filePath string, jsonFilePath string) (bool, error) {
	// Обчислення хешу файлу
	hash, err := calculateHash(filePath)
	if err != nil {
		return false, err
	}

	// Читання існуючих даних з JSON файлу
	var existingData []FileInfo
	if _, err := os.Stat(jsonFilePath); !os.IsNotExist(err) {
		file, err := os.Open(jsonFilePath)
		if err != nil {
			return false, err
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		err = decoder.Decode(&existingData)
		if err != nil && err != io.EOF {
			return false, err
		}
	}

	// Пошук інформації про файл в JSON
	for i, info := range existingData {
		if info.Path == filePath {
			// Хеш змінився, оновлюємо інформацію
			if info.Hash != hash {
				existingData[i].Hash = hash

				// Запис оновлених даних в JSON файл
				file, err := os.Create(jsonFilePath)
				if err != nil {
					return false, err
				}
				defer file.Close()

				encoder := json.NewEncoder(file)
				err = encoder.Encode(existingData)
				if err != nil {
					return false, err
				}

				return true, nil // Хеш змінився і був оновлений
			}
			return false, nil // Хеш не змінився
		}
	}

	// Файл не знайдено в JSON, додаємо новий запис
	existingData = append(existingData, FileInfo{
		Path: filePath,
		Name: filepath.Base(filePath),
		Hash: hash,
	})

	// Запис оновлених даних в JSON файл
	file, err := os.Create(jsonFilePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(existingData)
	if err != nil {
		return false, err
	}

	return true, nil // Файл додано або оновлено
}

// calculateHash обчислює SHA-256 хеш файлу
func calculateHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
