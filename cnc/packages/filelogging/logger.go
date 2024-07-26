package filelogging

import (
	"fmt"
	"os"
	"time"
)

type Logger struct {
	FilePath string
	file     *os.File
}

func NewLogger(path string) *Logger {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	return &Logger{
		FilePath: path,
		file:     f,
	}
}

func (logger *Logger) Logf(format string, val ...any) {
	_, _ = logger.file.WriteString("[" + time.Now().Format(time.RFC822) + "] " + fmt.Sprintf(format, val...) + "\r\n")
}

func (logger *Logger) Log(format string) {
	_, _ = logger.file.WriteString("[" + time.Now().Format(time.RFC822) + "] " + fmt.Sprint(format) + "\r\n")
}
