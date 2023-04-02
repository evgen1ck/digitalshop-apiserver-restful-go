package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
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

type Logger struct {
	*zap.Logger
}

func NewZap() (*Logger, error) {
	if err := os.MkdirAll(folderName, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	logFileName := fmt.Sprintf("%s.%s", time.Now().Format(fileFormat), fileExtension)
	logFilePath := filepath.Join(folderName, logFileName)

	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    64, // Maximum file size in megabytes
		MaxBackups: 3,  // Maximum number of rotated files to be saved
		MaxAge:     7,  // The maximum number of days during which the rotated files are stored
		Compress:   true,
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, zapcore.AddSync(lumberjackLogger), zapcore.InfoLevel)
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
