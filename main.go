package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sync"
)

func startTerminal(wg *sync.WaitGroup, input <-chan string, output chan<- string) {
	defer wg.Done()

	var cmd *exec.Cmd

	// Вибір терміналу залежно від операційної системи
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe")
	} else {
		cmd = exec.Command("bash")
	}

	// Створюємо канали для читання і запису
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	// Запускаємо команду (термінал)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	// Читаємо з терміналу (stdout)
	go func() {
		reader := bufio.NewReader(stdout)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Println("Error reading stdout:", err)
			}
			output <- line
		}
	}()

	// Читаємо помилки з терміналу (stderr)
	go func() {
		reader := bufio.NewReader(stderr)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Println("Error reading stderr:", err)
			}
			output <- line
		}
	}()

	// Читаємо команди від основного потоку
	for command := range input {
		if command == "exit" {
			stdin.Close()
			cmd.Process.Kill()
			close(output)
			return
		}
		_, err := io.WriteString(stdin, command+"\n")
		if err != nil {
			log.Println("Error writing to stdin:", err)
		}
	}

	// Чекаємо завершення процесу
	cmd.Wait()
}

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

	manager := management.Manager{}
	manager.Start()*/
	wg := &sync.WaitGroup{}
	input := make(chan string)
	output := make(chan string)

	wg.Add(1)
	go startTerminal(wg, input, output)

	// Основний цикл для взаємодії з користувачем
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter command: ")
		if scanner.Scan() {
			command := scanner.Text()

			if command == "exit" {
				input <- "exit"
				break
			}

			// Передаємо команду в потік термінала
			input <- command

			// Отримуємо та виводимо відповідь
			for {
				select {
				case line, ok := <-output:
					if !ok {
						break
					}
					fmt.Print(line)
				}
				break
			}
		}
	}

	close(input)
	wg.Wait()
	fmt.Println("Program terminated.")
}
