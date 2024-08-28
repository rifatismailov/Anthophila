package logging

import (
	"Anthophila/information"
	"encoding/json"
	"log"
	"net"
	"sync"
	"time"
)

// Logger - структура, яка використовується для зберігання повідомлень і помилок.
// Вона містить два поля: Message і Error, які використовуються для логування інформації.
type Logger struct {
	Message string `json:"Message"`
	Error   string `json:"Error"`
}

// Send відправляє лог до сервера логів у окремій горутині.
// Після завершення відправки горутина автоматично закривається.
//
// Параметри:
// - logServer: адреса сервера логів, до якого відправляється повідомлення.
// - wg: покажчик на `sync.WaitGroup`, який використовується для синхронізації завершення горутин.
// - done: канал, який сигналізує про завершення роботи горутини.
func (l *Logger) Send(logServer string, wg *sync.WaitGroup, done chan struct{}) {
	wg.Add(1) // Додаємо до лічильника горутин

	go func() {
		defer wg.Done() // Зменшуємо лічильник після завершення горутини

		// Внутрішня структура для відправки повідомлення
		type message struct {
			HostName    string `json:"HostName"`
			HostAddress string `json:"HostAddress"`
			MACAddress  string `json:"MACAddress"`
			RemoteAddr  string `json:"RemoteAddr"`
			Message     string `json:"Message"`
			Error       string `json:"Error"`
		}

		// Заповнення внутрішньої структури даними
		msg := message{
			HostName:    information.NewInfo().HostName(),
			HostAddress: information.NewInfo().HostAddress(),
			MACAddress:  information.NewInfo().GetMACAddress(),
			RemoteAddr:  information.NewInfo().RemoteAddress("https://api.ipify.org"),
			Message:     l.Message,
			Error:       l.Error,
		}

		// Серіалізація структури в JSON формат
		jsonData, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Error marshalling JSON: %v", err)
			return
		}

		// Відправка JSON даних до сервера логів
		sendLogger(logServer, string(jsonData))

		// Закриття горутини після успішного завершення відправки
		close(done)
	}()
}

// sendLogger встановлює з'єднання з сервером і відправляє логові дані у форматі JSON.
//
// Параметри:
// - serverAddress: адреса сервера, до якого слід підключитися.
// - json: серіалізовані JSON дані, які потрібно відправити.
func sendLogger(serverAddress, json string) {
	var conn net.Conn
	var err error

	for {
		// Спроба підключення до сервера
		conn, err = net.Dial("tcp", serverAddress)
		if err != nil {
			// Помилка підключення до сервера, спроба повторного підключення через 1 секунду
			time.Sleep(1 * time.Second)
			continue
		}
		// Підключення успішне, вихід з циклу перепідключення
		break
	}

	// Відправка повідомлення після успішного підключення
	defer conn.Close() // Закриття з'єднання після відправки
	_, err = conn.Write([]byte(json))
	if err != nil {
		log.Printf("Error sending message: %v", err)
	} else {
		log.Println("Message sent successfully")
	}
}
