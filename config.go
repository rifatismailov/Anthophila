package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Config структура для зберігання конфігураційних параметрів.
type Config struct {
	FileServer    string   `json:"file_server"`    // Адреса файлового сервера
	ManagerServer string   `json:"manager_server"` // Адреса сервера менеджера
	LogServer     string   `json:"log_server"`     // Адреса сервера логування
	Directories   []string `json:"directories"`    // Список директорій для перевірки
	Extensions    []string `json:"extensions"`     // Список розширень файлів для перевірки
	Hour          int      `json:"hour"`           // Година запуску
	Minute        int      `json:"minute"`         // Хвилина запуску
	Key           string   `json:"key"`            // Ключ для шифрування
}

const configFile = "config.json"

// loadConfig зчитує конфігураційний файл у форматі JSON. Якщо файл не існує або виникає помилка при читанні, повертає помилку.
func loadConfig() (*Config, error) {
	file, err := os.Open(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // Файл не існує
		}
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// saveConfig зберігає конфігурацію у файл у форматі JSON.
func saveConfig(config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configFile, data, 0644)
}
