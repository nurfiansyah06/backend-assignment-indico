package logs

import (
	"fmt"
	"log"
	"os"
	"time"
)

func InitLog() (*os.File, error) {
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, err
	}

	currentTime := time.Now()
	filename := fmt.Sprintf("logs/%s.log", currentTime.Format("01-02-2006"))

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	log.SetOutput(file)
	return file, nil
}

func Message(level string, msg string) {
	log.Printf("[%s] %s", level, msg)
}

func Info(msg string) {
	Message("INFO", msg)
}

func Error(msg string) {
	Message("ERROR", msg)
}
