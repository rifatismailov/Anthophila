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
*	Встановлює з'єднання з сервером, відправляє ім'я файлу (вирівняне до 256 байт),
*	обчислює та відправляє MD5 хеш-сумму файлу, шифрує файл та відправляє зашифрований файл на сервер.
 */

func (f *FILESender) SenderFile(serverAddr, filePath string, key []byte) error {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		logging.Now().PrintLog("[SenderFile] Помилка з'єднання", err.Error())
		return err
	}
	defer conn.Close()

	fileName := filepath.Base(filePath)
	fileNameBytes := []byte(fileName)

	// Перевірка довжини імені файлу
	if len(fileNameBytes) > 256 {
		logging.Now().PrintLog("[SenderFile] Ім'я файлу занадто довге", fileName)
		return err
	}

	// Вирівнювання імені файлу до 256 байт
	paddedFileNameBytes := make([]byte, 256)
	copy(paddedFileNameBytes, fileNameBytes)

	_, err = conn.Write(paddedFileNameBytes)
	if err != nil {
		logging.Now().PrintLog("[SenderFile] Помилка відправки імені файлу", err.Error())
		return err
	}

	// Відкриття файлу
	file, err := os.Open(filePath)
	if err != nil {
		logging.Now().PrintLog("[SenderFile] Помилка відкриття файлу",
			"{FilePath :{"+filePath+"} Err :{"+err.Error()+"}}")
		return err
	}
	defer file.Close()

	// Обчислення хеш-сумми файлу
	hasher := md5.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		logging.Now().PrintLog(
			"[SenderFile] Помилка обчислення хеш-сумми для файлу",
			"{FilePath :{"+filePath+"} Err :{"+err.Error()+"}}")
		return err
	}
	fileHash := hasher.Sum(nil)

	_, err = conn.Write(fileHash)
	if err != nil {
		logging.Now().PrintLog(
			"[SenderFile] Помилка відправки хеш-сумми файлу", err.Error())
		return err
	}

	// Зашифрування файлу та відряджання зашифрованого файлу на сервер
	file.Seek(0, 0)

	encrypt := cryptofile.NewFILEEncryptor()
	encryptedFile, err := encrypt.EncryptingFile(file, key)
	if err != nil {
		logging.Now().PrintLog("[SenderFile] Помилка шифрування файлу",
			"{FilePath :{"+filePath+"} Err :{"+err.Error()+"}}")
		return err
	}
	defer encryptedFile.Close()

	_, err = io.Copy(conn, encryptedFile)
	if err != nil {
		logging.Now().PrintLog(
			"[SenderFile] Помилка відправки зашифрованого файлу",
			"{FilePath :{"+filePath+"} Err :{"+err.Error()+"}}")
		return err
	}

	//printLog("[SenderFile] Зашифрований файл відправлено", filePath)
	err = deleteFile(filePath + ".enc")
	if err != nil {
		logging.Now().PrintLog(
			"Помилка при видаленні Зашифрованого файлу", err.Error())
	} else {
		//fmt.Println("Файл успішно видалено.")
	}
	return nil
}

// deleteFile видаляє файл за вказаним шляхом
func deleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}
