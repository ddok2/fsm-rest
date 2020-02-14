package utils

import (
	"io"
	"log"
	"os"
	"time"
)

var (
	Info    *log.Logger
	Trace   *log.Logger
	Warning *log.Logger
	Error   *log.Logger

	logger *Logger
)

var (
	colorOff    = "\033[0m"
	colorRed    = "\033[0;31m"
	colorGreen  = "\033[0;32m"
	colorOrange = "\033[0;33m"
	colorBlue   = "\033[0;34m"
	colorPurple = "\033[0;35m"
	colorCyan   = "\033[0;36m"
	colorGray   = "\033[0;37m"
)

type Logger struct {
	fpLog *os.File
}

func NewLogger() *Logger {
	if logger == nil {
		logger = new(Logger)
	}

	return logger
}

func (l *Logger) InitLogger() {
	var err error

	os.Mkdir("log", os.ModePerm)

	logFileName := "./log/automation_tester-" + time.Now().Format("2006-01-02-150405") + ".log"

	l.fpLog, err = os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	Info = log.New(io.MultiWriter(l.fpLog, os.Stdout), "INFO: ", log.Ldate|log.Lmicroseconds)
	Warning = log.New(io.MultiWriter(l.fpLog, os.Stdout), "WARNING: ", log.Ldate|log.Lmicroseconds)
	Trace = log.New(io.MultiWriter(l.fpLog, os.Stdout), "TRACE: ", log.Ldate|log.Lmicroseconds)
	Error = log.New(io.MultiWriter(l.fpLog, os.Stderr), "ERROR: ", log.Ldate|log.Lmicroseconds)

}

func (l *Logger) Info(v ...interface{}) {
	Info.Println(v)
}

func (l *Logger) Infof(template string, v ...interface{}) {
	Info.Printf(template, v)
}

func (l *Logger) Warn(v ...interface{}) {
	Info.Println(v)
}

func (l *Logger) Warnf(template string, v ...interface{}) {
	Info.Printf(template, v)
}

func (l *Logger) Trace(v ...interface{}) {
	Info.Println(v)
}

func (l *Logger) Tracef(template string, v ...interface{}) {
	Info.Printf(template, v)
}

func (l *Logger) Error(v ...interface{}) {
	Info.Println(v)
}

func (l *Logger) Errorf(template string, v ...interface{}) {
	Info.Printf(template, v...)
}

func (l *Logger) Finalize() {
	l.fpLog.Close()
}
