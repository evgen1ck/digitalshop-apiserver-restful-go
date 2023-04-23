package logger

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	tl "test-server-go/internal/tools"
	"time"
)

const (
	folderName    = "logs"
	fileFormat    = "02.01.06"
	fileExtension = "log"
)

var (
	logFileName = fmt.Sprintf("%s.%s", time.Now().Format(fileFormat), fileExtension)
)

type Logger struct {
	*zap.Logger
}

func NewZap() (*Logger, error) {
	var logFilePath string
	var err error

	logFilePath, err = tl.GetExecutablePathWithJoin(folderName, logFileName)
	if err != nil {
		return nil, err
	}

	if err = updateLogFile(&logFilePath); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, zapcore.AddSync(file), zapcore.InfoLevel)
	logger := zap.New(core)

	return &Logger{Logger: logger}, nil
}

func updateLogFile(logFilePath *string) error {
	today := time.Now().Format(fileFormat)
	currentLogFile := filepath.Base(*logFilePath)
	currentLogDate := currentLogFile[:len(today)]

	if currentLogDate != today {
		newLogFilePath, err := tl.GetExecutablePathWithJoin(folderName, logFileName)
		if err != nil {
			return err
		}

		*logFilePath = newLogFilePath
	}
	return nil
}

func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

func (l *Logger) NewInfo(message string) {
	l.Logger.Info(tl.CapitalizeFirst(message))
}

func (l *Logger) NewWarn(message string, error error) {
	if error == nil {
		error = errors.New("warn from zap logger")
	}

	_, file, line, _ := runtime.Caller(1)
	fields := []zapcore.Field{
		zap.String("file", file),
		zap.String("line", strconv.Itoa(line)),
	}

	l.Logger.Warn(tl.CapitalizeFirst(message)+": "+tl.UncapitalizeFirst(error.Error()), fields...)
}

func (l *Logger) NewError(message string, error error) {
	if error == nil {
		error = errors.New("error from zap logger")
	}

	_, file, line, _ := runtime.Caller(1)
	fields := []zapcore.Field{
		zap.String("file", file),
		zap.String("line", strconv.Itoa(line)),
	}

	l.Logger.Error(tl.CapitalizeFirst(message)+": "+tl.UncapitalizeFirst(error.Error()), fields...)
	os.Exit(1)
}
