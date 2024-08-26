package management

import (
	"github.com/gorilla/websocket"
	"log"
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
