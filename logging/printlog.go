package logging

type PrintLog struct {
}

func Now() *PrintLog {
	return &PrintLog{}
}
func (p PrintLog) PrintLog(message string, err string) {
	logger := Logger{Message: message, Error: err}
	logger.Send("localhost:6606")
}
