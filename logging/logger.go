package logging

import (
	"Anthophila/infirmation"
	"encoding/json"
	"net"
	"time"
)

type Logger struct {
	Message string `json:"Message"`
	Error   string `json:"Error"`
}

// Send відправляє лог в окремому потоці
func (l *Logger) Send(logServer string) {

	go func() {
		// Внутрішня структура для відряджання повідомлення
		type message struct {
			HostName    string `json:"HostName"`
			HostAddress string `json:"HostAddress"`
			MACAddress  string `json:"MACAddress"`
			RemoteAddr  string `json:"RemoteAddr"`
			Message     string `json:"Message"`
			Error       string `json:"Error"`
		}

		// Заповнення внутрішньої структури
		msg := message{
			HostName:    infirmation.NewInfo().HostName(),
			HostAddress: infirmation.NewInfo().HostAddress(),
			MACAddress:  infirmation.NewInfo().GetMACAddress(),
			RemoteAddr:  infirmation.NewInfo().RemoteAddress("https://api.ipify.org"),
			Message:     l.Message,
			Error:       l.Error,
		}

		// Серіалізація в JSON
		jsonData, err := json.Marshal(msg)
		if err != nil {
		}
		sendLogger(logServer, string(jsonData))
	}()
}

func sendLogger(serverAddress, json string) {
	var conn net.Conn
	var err error

	for {
		// Спроба підключення до сервера
		conn, err = net.Dial("tcp", serverAddress)
		if err != nil {
			//Помилка підключення до сервера
			time.Sleep(1 * time.Second) // Затримка перед спробою повторного підключення
			continue
		}
		// Підключення успішне, вихід з циклу перепідключення
		break
	}

	// Відряджання повідомлення після успішного підключення
	defer conn.Close() // Закриття з'єднання після відряджання
	_, err = conn.Write([]byte(json))
	if err != nil {
		//Помилка при надсиланні повідомлення
	} else {
		//Повідомлення надіслано успішно.
	}
}
