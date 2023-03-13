package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
)

// Logger is a struct that embeds logrus.Logger
type Logger struct {
	*logrus.Logger
}

// New creates a new Logger instance with a logrus logger and sets its output to os.Stdout
func New() *Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})
	return &Logger{logger}
}

// NewInfo logs an informational message with the given message
func (l *Logger) NewInfo(msg string) {
	l.WithFields(logrus.Fields{
		"severity": "info",
	}).Info(msg)
}

// NewWarn logs a warning message with the file name and line number of the calling function
func (l *Logger) NewWarn(msg string) {
	// Get the file name and line number of the calling function
	_, file, line, _ := runtime.Caller(1)

	l.WithFields(logrus.Fields{
		"severity": "warning",
		"file":     file,
		"line":     line,
	}).Warn(msg)
}

// NewErrorWithExit logs an error message with the file name and line number of the calling function and exits the program
func (l *Logger) NewErrorWithExit(msg string, err error) {
	_, file, line, _ := runtime.Caller(1)

	l.WithFields(logrus.Fields{
		"severity": "error",
		"file":     file,
		"line":     line,
	}).Error(msg + ": " + err.Error())
	os.Exit(1)
}

// NewError logs an error message with the file name and line number of the calling function
func (l *Logger) NewError(msg string, err error) {
	_, file, line, _ := runtime.Caller(1)

	l.WithFields(logrus.Fields{
		"severity": "error",
		"file":     file,
		"line":     line,
	}).Error(msg + ": " + err.Error())
}
