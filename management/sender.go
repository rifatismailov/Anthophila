package management

import (
	"bufio"
	"github.com/gorilla/websocket"
	"log"
	"os"
)

type Sender struct{}

func NewSender() *Sender {
	return &Sender{}
}
func (*Sender) sendMessage(wSocket *websocket.Conn) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		log.Print("Enter message: ")
		scanner.Scan()
		text := scanner.Text()

		if text == "exit" {
			break
		}

		err := wSocket.WriteMessage(websocket.TextMessage, []byte(text))
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}
	}
}
func (*Sender) sendMessageWith(wSocket *websocket.Conn, text []byte) {
	err := wSocket.WriteMessage(websocket.TextMessage, text)
	if err != nil {
		log.Println("Error sending message:", err)
		return
	}
}
