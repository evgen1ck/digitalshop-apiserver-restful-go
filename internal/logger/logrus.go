package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
)

type Logger struct {
	*logrus.Logger
}

func New() *Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})
	return &Logger{logger}
}

func (l *Logger) NewInfo(msg string) {
	l.WithFields(logrus.Fields{
		"severity": "info",
	}).Info(msg)
}

func (l *Logger) NewWarn(msg string) {
	// Get the file name and line number of the calling function
	_, file, line, _ := runtime.Caller(1)

	l.WithFields(logrus.Fields{
		"severity": "warning",
		"file":     file,
		"line":     line,
	}).Warn(msg)
}

func (l *Logger) NewError(msg string, err error) {
	_, file, line, _ := runtime.Caller(1)

	l.WithFields(logrus.Fields{
		"severity": "error",
		"file":     file,
		"line":     line,
	}).Error(msg + ": " + err.Error())
	os.Exit(1)
}
