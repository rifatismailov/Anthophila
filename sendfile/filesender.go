package sendfile

import (
	"Anthophila/cryptofile"
	"crypto/md5"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
)

type FILESender struct {
	// ... (поля для зберігання стану, якщо потрібні)
}

func NewFILESender() *FILESender {
	return &FILESender{}
}

/*
*	Функція для відправки файлу на сервер sendFileToServer:
*	Встановлює з'єднання з сервером, відправляє ім'я файлу (вирівняне до 256 байт),
*	обчислює та відправляє MD5 хеш-сумму файлу, шифрує файл та відправляє зашифрований файл на сервер.
 */
func (f *FILESender) SenderFile(serverAddr, filePath string, key []byte) error {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Printf("Помилка з'єднання: %v\n", err)
		return err
	}
	defer conn.Close()

	fileName := filepath.Base(filePath)
	fileNameBytes := []byte(fileName)

	// Перевірка довжини імені файлу
	if len(fileNameBytes) > 256 {
		fmt.Printf("Ім'я файлу занадто довге: %s\n", fileName)
		return err
	}

	// Вирівнювання імені файлу до 256 байт
	paddedFileNameBytes := make([]byte, 256)
	copy(paddedFileNameBytes, fileNameBytes)

	_, err = conn.Write(paddedFileNameBytes)
	if err != nil {
		fmt.Printf("Помилка відправки імені файлу: %v\n", err)
		return err
	}

	// Відкриття файлу
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Помилка відкриття файлу %s: %v\n", filePath, err)
		return err
	}
	defer file.Close()

	// Обчислення хеш-сумми файлу
	hasher := md5.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		fmt.Printf("Помилка обчислення хеш-сумми для файлу %s: %v\n", filePath, err)
		return err
	}
	fileHash := hasher.Sum(nil)

	_, err = conn.Write(fileHash)
	if err != nil {
		fmt.Printf("Помилка відправки хеш-сумми файлу: %v\n", err)
		return err
	}

	// Зашифрування файлу та відправка зашифрованого файлу на сервер
	file.Seek(0, 0)

	encrypt := cryptofile.NewFILEEncryptor()
	encryptedFile, err := encrypt.EncryptingFile(file, key)
	if err != nil {
		fmt.Printf("Помилка шифрування файлу %s: %v\n", filePath, err)
		return err
	}
	defer encryptedFile.Close()

	_, err = io.Copy(conn, encryptedFile)
	if err != nil {
		fmt.Printf("Помилка відправки зашифрованого файлу %s: %v\n", filePath, err)
		return err
	}

	fmt.Printf("Зашифрований файл відправлено: %s\n", filePath)
	return nil
}
