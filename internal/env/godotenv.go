package env

import (
	"errors"
	"os"
	"strconv"
	"test-server-go/internal/logger"
)

func GetEnv(key string, logger *logger.Logger) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	logger.NewError("GetEnv", errors.New("not found "+key))
	return ""
}

func GetEnvAsBool(key string, logger *logger.Logger) bool {
	valStr := GetEnv(key, logger)
	if value, err := strconv.ParseBool(valStr); err == nil {
		return value
	}

	logger.NewError("GetEnvAsBool", errors.New("not convert "+key))
	return false
}

func GetEnvAsInt(key string, logger *logger.Logger) int {
	valStr := GetEnv(key, logger)
	if value, err := strconv.Atoi(valStr); err == nil {
		return value
	}

	logger.NewError("GetEnvAsInt", errors.New("not convert "+key))
	return 0
}
