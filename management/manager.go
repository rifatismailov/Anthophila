package management

import (
	"Anthophila/information"
	"bufio"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"os/signal"
	"time"
)

type Manager struct {
}

func (f *Manager) Start() {

	for {
		//{"sClient":"Bob","rClient":"Alex","message":"У JSON-повідомленні, яке ви надсилаєте, використовується неправильний регістр для цього ключа!"}

		// Підключення до сервера
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)

		// Підключення до WebSocket-сервера
		c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
		if err != nil {
			log.Fatal("dial:", err)
		}
		defer c.Close()

		done := make(chan struct{})

		go func() {
			defer close(done)
			for {
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Println("Error reading from server:", err)
					return
				}
				log.Printf("Received: %s", message)
			}
		}()
		macAddress := information.NewInfo().GetMACAddress()
		// Надсилання повідомлення
		err = c.WriteMessage(websocket.TextMessage, []byte("nick:"+macAddress))
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}
		scanner := bufio.NewScanner(os.Stdin)

		for {
			log.Print("Enter message: ")
			scanner.Scan()
			text := scanner.Text()

			if text == "exit" {
				break
			}

			err = c.WriteMessage(websocket.TextMessage, []byte(text))
			if err != nil {
				log.Println("Error sending message:", err)
				return
			}
		}

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				err := c.WriteMessage(websocket.TextMessage, []byte("Ping"))
				if err != nil {
					log.Println("Error writing to server:", err)
					return
				}
			case <-interrupt:
				log.Println("interrupt")

				err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Println("Error during close handshake:", err)
					return
				}
				select {
				case <-done:
				case <-time.After(time.Second):
				}
				return
			}
		}
	}

}
