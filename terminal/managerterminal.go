package terminal

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

// TManager представляє структуру для керування терміналом
type TManager struct {
	cmd    *exec.Cmd
	input  chan string
	output chan string
	wg     *sync.WaitGroup
	pid    int
	mu     sync.Mutex // Для синхронізації доступу до cmd
}

// NewTerminalManager створює новий екземпляр TManager
func NewTerminalManager() *TManager {
	tm := &TManager{
		input:  make(chan string),
		output: make(chan string),
		wg:     &sync.WaitGroup{},
	}

	// Вибір терміналу залежно від операційної системи
	if runtime.GOOS == "windows" {
		tm.cmd = exec.Command("cmd.exe")
	} else {
		tm.cmd = exec.Command("bash")
	}

	return tm
}

// Start запускає термінал та потоки для взаємодії з ним
func (tm *TManager) Start() error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Якщо термінал уже запущений, повертаємо помилку
	if tm.cmd != nil && tm.cmd.Process != nil && tm.cmd.ProcessState == nil {
		return errors.New("terminal already running")
	}

	// Оновлюємо команду для терміналу
	tm.cmd = exec.Command(tm.cmd.Path)

	stdin, err := tm.cmd.StdinPipe()
	if err != nil {
		log.Println("Error creating stdin pipe:", err)
		return err
	}
	stdout, err := tm.cmd.StdoutPipe()
	if err != nil {
		log.Println("Error creating stdout pipe:", err)
		return err
	}
	stderr, err := tm.cmd.StderrPipe()
	if err != nil {
		log.Println("Error creating stderr pipe:", err)
		return err
	}

	if err := tm.cmd.Start(); err != nil {
		log.Println("Error starting command:", err)
		return err
	}

	tm.pid = tm.cmd.Process.Pid

	tm.wg.Add(1)
	go tm.runTerminal(stdin, stdout, stderr)

	return nil
}

// Stop зупиняє термінал
func (tm *TManager) Stop() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if tm.cmd.Process == nil {
		return
	}

	tm.input <- "exit"
	tm.cmd.Process.Kill() // Завершуємо процес
	tm.cmd = nil          // Скидаємо cmd
	tm.wg.Wait()
}

// SendCommand надсилає команду до терміналу
func (tm *TManager) SendCommand(command string) {
	tm.input <- command
}

// GetOutput повертає канал для читання виходу терміналу
func (tm *TManager) GetOutput() <-chan string {
	return tm.output
}

// Restart перезапускає термінал
func (tm *TManager) Restart() {
	tm.Stop() // Зупиняємо термінал
	time.Sleep(1 * time.Second)
	if err := tm.Start(); err != nil {
		log.Println("Failed to start terminal:", err)
	}
}

// runTerminal запускає обробку вводу/виводу термінала
func (tm *TManager) runTerminal(stdin io.WriteCloser, stdout io.Reader, stderr io.Reader) {
	defer tm.wg.Done()

	go func() {
		defer close(tm.output)
		reader := bufio.NewReader(stdout)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Println("Error reading stdout:", err)
				break
			}
			tm.output <- line
		}
	}()

	go func() {
		reader := bufio.NewReader(stderr)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Println("Error reading stderr:", err)
				break
			}
			tm.output <- line
		}
	}()

	for command := range tm.input {
		if strings.TrimSpace(command) == "exit" {
			stdin.Close()
			return
		}

		if strings.TrimSpace(command) == "stop" {
			continue
		}

		if strings.TrimSpace(command) == "sudo su" {
			fmt.Println("YOU MUST BE SUDO:")
			continue
		}

		if strings.HasPrefix(command, "ping ") && !strings.Contains(command, "-c") {
			// Заміна команди на ping -c 4 example.com
			parts := strings.Split(command, " ")
			if len(parts) > 1 {
				command = fmt.Sprintf("ping -c 4 %s", strings.Join(parts[1:], " "))
			}
		}

		_, err := io.WriteString(stdin, command+"\n")
		if err != nil {
			log.Println("Error writing to stdin:", err)
			tm.Restart() // Перезапуск термінала при помилці
		}
	}

	if err := tm.cmd.Wait(); err != nil {
		log.Println("Error waiting for command:", err)
	}
}
