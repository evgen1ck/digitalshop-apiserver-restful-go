package config

import (
	"errors"
	"os"
	"strconv"
	"test-server-go/internal/logger"
)

func getEnv(key string, logger *logger.Logger) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	logger.NewErrorWithExit("getEnv", errors.New("not found "+key))
	return ""
}

func getEnvAsBool(key string, logger *logger.Logger) bool {
	valStr := getEnv(key, logger)
	if value, err := strconv.ParseBool(valStr); err == nil {
		return value
	}

	logger.NewErrorWithExit("getEnvAsBool", errors.New("not convert "+key))
	return false
}

func getEnvAsInt(key string, logger *logger.Logger) int {
	valStr := getEnv(key, logger)
	if value, err := strconv.Atoi(valStr); err == nil {
		return value
	}

	logger.NewErrorWithExit("getEnvAsInt", errors.New("not convert "+key))
	return 0
}
