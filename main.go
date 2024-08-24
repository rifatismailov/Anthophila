package main

import (
	"Anthophila/checkfile"
	"fmt"
	"time"
)

func main() {
	directories := []string{"/home/sirius/GolandProjects/Anthophila/Documents/"}
	file_cheker := checkfile.FileChecker{Address: "localhost:12345",
		Key: []byte("a very very very very secret key"), Directories: directories,
		SupportedExtensions: []string{".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx"}}
	file_cheker.Start()
	// Головний потік може продовжувати свою роботу, наприклад, обробляти інші завдання
	fmt.Println("Main goroutine continues...")

	for {
		// Цей цикл буде виконуватися вічно
		time.Sleep(time.Second) // Можна змінити затримку за потреби
	}
}
