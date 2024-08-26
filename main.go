package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	/*infoJson := information.NewInfo().InfoJson()
	directories := []string{"/home/sirius/GolandProjects/Anthophila/"}
	file_cheker := checkfile.FileChecker{
		Address:             "localhost:12345",
		Key:                 []byte("a very very very very secret key"),
		Directories:         directories,
		SupportedExtensions: []string{".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx"},
		TimeStart:           []int8{10, 25},
		InfoJson:            infoJson}
	file_cheker.Start()
	for {

		// Цей цикл буде виконуватися вічно
		fmt.Println("Main goroutine continues...")
		time.Sleep(time.Second) // Можна змінити затримку за потреби
	}
	*/

	// Підключення до сервера
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server...")

	// Запуск горутини для отримання повідомлень від сервера
	go func() {
		reader := bufio.NewReader(conn)
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading from server:", err)
				return
			}
			fmt.Print("Server: " + message)
			//{"sClient":"B","rClient":"A","Message":"message"}
		}
	}()

	// Основний цикл для введення повідомлень
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter message: ")
		if scanner.Scan() {
			message := strings.TrimSpace(scanner.Text())

			// Перевірка на порожній рядок
			if len(message) == 0 {
				fmt.Println("Empty message. Please enter a valid command or message.")
				continue
			}

			// Відправлення повідомлення на сервер
			_, err := conn.Write([]byte(message + "\n"))
			if err != nil {
				fmt.Println("Error sending message:", err)
				return
			}
		}
	}
}
