package management

import (
	"Anthophila/information"
	"Anthophila/logging"
	"github.com/gorilla/websocket"
	"time"
)

const (
	reconnectInterval = 5 * time.Second // Інтервал для повторних спроб підключення
)

// Manager структура, що представляє управління підключенням до WebSocket сервера
type Manager struct {
}

// Start ініціалізує підключення до WebSocket сервера і обробляє повідомлення
func (f *Manager) Start(logStatus bool, logAddress, serverAddr string) {
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
			if logStatus == true {
				logging.Now().PrintLog(logAddress, "Error connecting to server: %v", err.Error())
				logging.Now().PrintLog(logAddress, "Retrying in %v...", reconnectInterval.String())
			}
			// Затримка перед наступною спробою підключення
			time.Sleep(reconnectInterval)
			continue
		}

		// Підключення успішне, надсилаємо нікнейм
		err = wSocket.WriteMessage(websocket.TextMessage, []byte("nick:"+macAddress))
		if err != nil {
			// Логування помилки відправки нікнейму
			if logStatus == true {
				logging.Now().PrintLog(logAddress, "Error sending nickname: %v", err.Error())
			}
			wSocket.Close()
			if logStatus == true {
				logging.Now().PrintLog(logAddress, "Retrying in %v...", reconnectInterval.String())
			}
			// Затримка перед наступною спробою підключення
			time.Sleep(reconnectInterval)
			continue
		}

		// Запуск горутіни для обробки отриманих повідомлень від сервера
		go NewReader().ReadMessageCommand(logStatus, logAddress, wSocket)

		// Основний цикл для надсилання повідомлень до сервера
		for {
			select {
			case <-time.After(reconnectInterval):
				// Надсилання пінгу для перевірки з'єднання
				errPing = wSocket.WriteMessage(websocket.TextMessage, []byte("Ping"))
			}

			if errPing != nil {
				// Логування помилки при надсиланні пінгу
				if logStatus == true {
					logging.Now().PrintLog(logAddress, "Error writing to server: %v", errPing.Error())
				}
				// Вихід з циклу при помилці
				break
			}
		}

		// Якщо ми потрапили сюди, це означає, що з'єднання було розірвано
		if logStatus == true {
			logging.Now().PrintLog(logAddress, "Connection closed. Retrying in %v...", reconnectInterval.String())
		}
		wSocket.Close()
		// Затримка перед наступною спробою підключення
		time.Sleep(reconnectInterval)
	}
}
