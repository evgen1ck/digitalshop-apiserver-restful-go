package logger

import (
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

type Logger struct {
	*zap.Logger
}

func NewZap(format string) (*Logger, error) {
	logPath := "logs"
	if err := os.MkdirAll(logPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	logFileName := fmt.Sprintf("%s", time.Now().Format("02.01.06"))
	if format == "gz" {
		logFileName += ".log.gz"
	} else {
		logFileName += ".log.zip"
	}

	logFile, err := os.OpenFile(filepath.Join(logPath, logFileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, zapcore.AddSync(logFile), zapcore.InfoLevel)
	logger := zap.New(core)

	return &Logger{Logger: logger}, nil
}

func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

func (l *Logger) NewInfo(message string) {
	l.Logger.Info(tl.CapitalizeFirst(message))
}

func (l *Logger) NewWarn(message string, error error) {
	_, file, line, _ := runtime.Caller(1)
	fields := []zapcore.Field{
		zap.String("file", file),
		zap.String("line", strconv.Itoa(line)),
	}

	l.Logger.Warn(tl.CapitalizeFirst(message)+": "+tl.UncapitalizeFirst(error.Error()), fields...)
}

func (l *Logger) NewError(message string, error error) {
	_, file, line, _ := runtime.Caller(1)
	fields := []zapcore.Field{
		zap.String("file", file),
		zap.String("line", strconv.Itoa(line)),
	}

	l.Logger.Error(tl.CapitalizeFirst(message)+": "+tl.UncapitalizeFirst(error.Error()), fields...)
	os.Exit(1)
}
