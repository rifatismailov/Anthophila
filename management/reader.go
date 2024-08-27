package management

import (
	"Anthophila/information"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

type Reader struct{}

func NewReader() *Reader {
	return new(Reader)
}

//{"sClient":"Alex","rClient":"Bob","message":"У JSON-повідомленні, яке ви надсилаєте, використовується неправильний регістр для цього ключа!"}

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

// Обробка отриманих повідомлень
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

	for {
		_, message, err := wSocket.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}

		eRR := json.Unmarshal(message, &cmd)
		if eRR != nil {
			// Якщо розпарсити як JSON не вдалося, обробляємо як звичайний текст
			log.Println("Received text message:", string(message))
		} else {
			// Якщо розпарсити вдалося, працюємо з даними
			//якщо команда дорівнює "help"
			if cmd.Command == "help" {
				// Заповнення внутрішньої структури
				msg := myMessage{
					SClient: information.NewInfo().GetMACAddress(),
					RClient: cmd.SClient,
					Message: "Ми отримали від вас повідомлення Наше імя : " + information.NewInfo().GetMACAddress() + " ваще повідомлення" + string(message),
				}

				// Серіалізація в JSON
				jsonData, err := json.Marshal(msg)
				if err != nil {
				}
				log.Println("Json " + string(jsonData))
				go NewSender().sendMessageWith(wSocket, jsonData)
			}
		}

	}
}
