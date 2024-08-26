package management

import (
	"Anthophila/information"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	serverAddr        = "ws://localhost:8080/ws"
	reconnectInterval = 5 * time.Second
)

type Manager struct {
}

func (f *Manager) Start() {
	macAddress := information.NewInfo().GetMACAddress()

	var wSocket *websocket.Conn
	var err error
	var errPing error

	for {
		// Спроба підключення до сервера
		wSocket, _, err = websocket.DefaultDialer.Dial(serverAddr, nil)
		if err != nil {
			log.Printf("Error connecting to server: %v", err)
			log.Printf("Retrying in %v...", reconnectInterval)
			time.Sleep(reconnectInterval)
			continue
		}

		// Підключення успішне, надсилаємо нікнейм
		err = wSocket.WriteMessage(websocket.TextMessage, []byte("nick:"+macAddress))
		if err != nil {
			log.Printf("Error sending nickname: %v", err)
			wSocket.Close()
			log.Printf("Retrying in %v...", reconnectInterval)
			time.Sleep(reconnectInterval)
			continue
		}

		// Запуск горутіни для отримання повідомлень від сервера
		go NewReader().ReadMessage(wSocket)
		// Запуск горутіни для відправки повідомлень від сервера
		go NewSender().sendMessage(wSocket)
		// Основний цикл для надсилання повідомлень
		for {
			select {
			case <-time.After(reconnectInterval):
				errPing = wSocket.WriteMessage(websocket.TextMessage, []byte("Ping"))
			}
			if errPing != nil {
				log.Printf("Error writing to server: %v", errPing)

				break
			}
		}

		// Якщо ми потрапили сюди, це означає, що з'єднання було розірвано
		log.Printf("Connection closed. Retrying in %v...", reconnectInterval)
		wSocket.Close()
		time.Sleep(reconnectInterval)
	}
}
