package logging

import (
	"log"
	"os"
)

type logWriter struct{}

func (writer logWriter) Write(p []byte) (n int, err error) {
	file, err := os.OpenFile("backend.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file: ", err)
	}
	_, err1 := file.Write(p)
	if err1 != nil {
		return 0, err1
	}
	_, err2 := os.Stdout.Write(p)
	if err2 != nil {
		return 0, err2
	}
	return len(p), nil
}
func ConfigureLogger() {
	log.SetOutput(logWriter{})
}
