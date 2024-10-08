package logging

import (
	"sync"
)

// PrintLog визначає структуру для логування повідомлень.
// У цьому прикладі структура не містить полів, але може бути розширена для зберігання додаткових даних.
type PrintLog struct {
}

// Now створює та повертає новий екземпляр структури PrintLog.
// Цей метод може бути корисний для ініціалізації або підготовки до логування.
func Now() *PrintLog {
	return &PrintLog{}
}

// PrintLog надсилає повідомлення та помилки до сервера логів.
// Він запускає горутину для відправки даних, а потім чекає, поки всі горутини завершаться.
//
// Параметри:
// - message: повідомлення, яке потрібно записати в лог.
// - err: повідомлення про помилку, яке потрібно записати в лог.
//
// Процес:
// 1. Ініціалізує `sync.WaitGroup` для відстеження горутин.
// 2. Створює канал `done` для сигналізації завершення горутини.
// 3. Створює екземпляр `Logger` з переданим повідомленням та помилкою.
// 4. Викликає метод `Send`, який запускає горутину для відправки даних до сервера логів.
// 5. Використовує `wg.Wait()` для очікування завершення горутини.
func (p PrintLog) PrintLog(LogAddress string, message string, err string) {
	var wg sync.WaitGroup
	done := make(chan struct{})
	logger := Logger{Message: message, Error: err}
	logger.Send(LogAddress, &wg, done)

	// Очікуємо завершення горутин
	wg.Wait()
}
