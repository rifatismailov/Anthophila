package cryptofile

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
)

type FILEEncryptor struct {
	// ... (поля для зберігання стану, якщо потрібні)
}

func NewFILEEncryptor() *FILEEncryptor {
	return &FILEEncryptor{}
}

/*
 *	Функція EncryptingFile для шифрування файлу:
 *	Шифрує вміст файлу за допомогою AES-256 та повертає зашифрований файл.
 */

func (f *FILEEncryptor) EncryptingFile(file *os.File, key []byte) (*os.File, error) {
	// Додаємо суфікс ".enc" до імені оригінального файлу, щоб створити шлях для зашифрованого файлу
	encryptedFilePath := file.Name() + ".enc"

	// Створюємо новий файл для збереження зашифрованих даних
	encryptedFile, err := os.Create(encryptedFilePath)
	if err != nil {
		return nil, err // Якщо створити файл не вдалося, повертаємо помилку
	}

	// Створюємо AES-блок шифрування з переданим ключем
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err // Повертаємо помилку у разі невдачі
	}

	// Генеруємо вектор ініціалізації (IV) для шифрування
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err // Якщо генерація IV не вдалася, повертаємо помилку
	}

	// Записуємо IV у початок зашифрованого файлу
	_, err = encryptedFile.Write(iv)
	if err != nil {
		return nil, err // Якщо запис IV не вдався, повертаємо помилку
	}

	// Створюємо новий шифрувальний потік на основі блоку AES та вектора IV
	stream := cipher.NewCFBEncrypter(block, iv)
	writer := &cipher.StreamWriter{S: stream, W: encryptedFile}

	// Шифруємо та копіюємо дані з оригінального файлу у новий файл
	if _, err := io.Copy(writer, file); err != nil {
		return nil, err // Якщо копіювання даних не вдалося, повертаємо помилку
	}

	// Повертаємо вказівник на початок зашифрованого файлу
	encryptedFile.Seek(0, 0)

	// Повертаємо зашифрований файл
	return encryptedFile, nil
}
