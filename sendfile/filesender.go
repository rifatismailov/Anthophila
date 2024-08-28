package sendfile

import (
	"Anthophila/cryptofile"
	"Anthophila/logging"
	"crypto/md5"
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
*	Функція для відряджання файлу на сервер sendFileToServer:
*	Встановлює з'єднання з сервером, відправляє ім'я файлу (вирівняне до 512 байт),
*	обчислює та відправляє MD5 хеш-сумму файлу, шифрує файл та відправляє зашифрований файл на сервер.
 */

func (f *FILESender) SenderFile(fileAddress, logAddress, filePath string, key []byte, infoJson string) {

	conn, err := net.Dial("tcp", fileAddress)
	if err != nil {
		logging.Now().PrintLog(logAddress, "[SenderFile] Помилка з'єднання", err.Error())
	}
	defer conn.Close()

	// Modify the filename by adding a prefix or suffix
	fileName := filepath.Base(filePath)
	modifiedFileName := infoJson + fileName // Example: Add "encrypted_" as a prefix
	fileNameBytes := []byte(modifiedFileName)

	// Ensure the filename is not too long
	if len(fileNameBytes) > 512 {
		logging.Now().PrintLog(logAddress, "[SenderFile] Ім'я файлу занадто довге", modifiedFileName)
	}

	// Pad the filename to 512 bytes
	paddedFileNameBytes := make([]byte, 512)
	copy(paddedFileNameBytes, fileNameBytes)

	_, err = conn.Write(paddedFileNameBytes)
	if err != nil {
		logging.Now().PrintLog(logAddress, "[SenderFile] Помилка відправки імені файлу", err.Error())
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		logging.Now().PrintLog(logAddress, "[SenderFile] Помилка відкриття файлу",
			"{FilePath :{"+filePath+"} Err :{"+err.Error()+"}}")
	}
	defer file.Close()

	// Calculate the file's hash sum
	hasher := md5.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		logging.Now().PrintLog(logAddress, "[SenderFile] Помилка обчислення хеш-сумми для файлу",
			"{FilePath :{"+filePath+"} Err :{"+err.Error()+"}}")
	}
	fileHash := hasher.Sum(nil)

	_, err = conn.Write(fileHash)
	if err != nil {
		logging.Now().PrintLog(logAddress, "[SenderFile] Помилка відправки хеш-сумми файлу", err.Error())
	}

	// Encrypt the file and send it to the server
	file.Seek(0, 0)

	encrypt := cryptofile.NewFILEEncryptor()
	encryptedFile, err := encrypt.EncryptingFile(file, key)
	if err != nil {
		logging.Now().PrintLog(logAddress, "[SenderFile] Помилка шифрування файлу",
			"{FilePath :{"+filePath+"} Err :{"+err.Error()+"}}")
	}
	defer encryptedFile.Close()

	_, err = io.Copy(conn, encryptedFile)
	if err != nil {
		logging.Now().PrintLog(logAddress, "[SenderFile] Помилка відправки зашифрованого файлу",
			"{FilePath :{"+filePath+"} Err :{"+err.Error()+"}}")
	}

	// Delete the encrypted file locally
	err = deleteFile(encryptedFile.Name())
	if err != nil {
		logging.Now().PrintLog(logAddress, "Помилка при видаленні Зашифрованого файлу", err.Error())
	} else {
		//fmt.Println("Файл успішно видалено.")
	}
}

// deleteFile видаляє файл за вказаним шляхом
func deleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}
