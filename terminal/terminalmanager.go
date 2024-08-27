package terminal

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

// TerminalManager представляє структуру для керування терміналом
type TerminalManager struct {
	cmd    *exec.Cmd
	input  chan string
	output chan string
	wg     *sync.WaitGroup
	pid    int
}

// NewTerminalManager створює новий екземпляр TerminalManager
func NewTerminalManager() *TerminalManager {
	tm := &TerminalManager{
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
func (tm *TerminalManager) Start() {
	tm.wg.Add(1)
	go tm.runTerminal()
}

// Stop зупиняє термінал
func (tm *TerminalManager) Stop() {
	tm.input <- "exit"
	close(tm.input)
	tm.wg.Wait()
	close(tm.output)
}

// SendCommand надсилає команду до терміналу
func (tm *TerminalManager) SendCommand(command string) {
	tm.input <- command
}

// GetOutput повертає канал для читання виходу терміналу
func (tm *TerminalManager) GetOutput() <-chan string {
	return tm.output
}
func (tm *TerminalManager) Restart() {
	tm.Stop()  // Зупиняємо горутину
	tm.Start() // Перезапускаємо горутину
}

// Внутрішній метод для керування терміналом
func (tm *TerminalManager) runTerminal() {

	defer tm.wg.Done()

	stdin, err := tm.cmd.StdinPipe()
	if err != nil {
		log.Fatal("Error creating stdin pipe:", err)
	}
	stdout, err := tm.cmd.StdoutPipe()
	if err != nil {
		log.Fatal("Error creating stdout pipe:", err)
	}
	stderr, err := tm.cmd.StderrPipe()
	if err != nil {
		log.Fatal("Error creating stderr pipe:", err)
	}

	if err := tm.cmd.Start(); err != nil {
		log.Fatal("Error starting command:", err)
	}

	tm.pid = tm.cmd.Process.Pid // Зберігаємо PID

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
			tm.cmd.Process.Kill()
			return
		}

		if strings.TrimSpace(command) == "stop" {
			// Завершити процес за PID

			continue
		}
		if strings.TrimSpace(command) == "sudo su" {
			fmt.Println("YOU MAST BE SUDO:")
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
			tm.Start()
		}
	}

	if err := tm.cmd.Wait(); err != nil {
		log.Println("Error waiting for command:", err)
	}
}

// Надсилає сигнал до процесу
func (tm *TerminalManager) sendSignalToProcess(signal string) {
	log.Println("Kill Procces:")
	var signalCmd *exec.Cmd

	signalCmd = exec.Command("kill", signal, fmt.Sprintf("%d", tm.cmd.Process.Pid))

	if err := signalCmd.Run(); err != nil {
		log.Println("Error sending signal to process:", err)
	}
}
