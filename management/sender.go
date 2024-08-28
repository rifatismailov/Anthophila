package management

import (
	"Anthophila/logging"
	"github.com/gorilla/websocket"
)

type Sender struct{}

func NewSender() *Sender {
	return &Sender{}
}

func (*Sender) sendMessageWith(logAddress string, wSocket *websocket.Conn, text []byte) error {
	err := wSocket.WriteMessage(websocket.TextMessage, text)
	if err != nil {
		logging.Now().PrintLog(logAddress, "Error sending message:", err.Error())
		return err
	}
	return err
}
