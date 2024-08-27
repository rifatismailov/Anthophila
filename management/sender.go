package management

import (
	"github.com/gorilla/websocket"
	"log"
)

type Sender struct{}

func NewSender() *Sender {
	return &Sender{}
}

func (*Sender) sendMessageWith(wSocket *websocket.Conn, text []byte) error {
	err := wSocket.WriteMessage(websocket.TextMessage, text)
	if err != nil {
		log.Println("Error sending message:", err)
		return err
	}
	return err
}
