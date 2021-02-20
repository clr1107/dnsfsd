package logger

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

func (l *Logger) Init(path string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		return err
	}

	writer := io.MultiWriter(os.Stdout, file)
	l.infoLogger = log.New(writer, "INFO: ", log.Ldate|log.Ltime)
	l.errorLogger = log.New(writer, "ERROR: ", log.Ldate|log.Ltime)

	return nil
}

func (l *Logger) log0(err bool, msg string, v ...interface{}) {
	x := l.infoLogger

	if err {
		x = l.errorLogger
	}

	if msg[len(msg)-1] != '\n' {
		msg += "\n"
	}

	x.Printf(msg, v...)
}

func (l *Logger) Log(msg string, v ...interface{}) {
	l.log0(false, msg, v...)
}

func (l *Logger) LogErr(msg string, v ...interface{}) {
	l.log0(true, msg, v...)
}

func (l *Logger) LogFatal(msg string, v ...interface{}) {
	l.LogErr(msg, v...)
	os.Exit(1)
}
