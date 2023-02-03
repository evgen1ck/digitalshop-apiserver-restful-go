package env

import (
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func GetEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	log.Fatal(errors.New("GetEnv: not found " + key))
	return ""
}

func GetEnvAsBool(key string) bool {
	valStr := GetEnv(key)
	if value, err := strconv.ParseBool(valStr); err == nil {
		return value
	}

	log.Fatal(errors.New("GetEnvAsBool: not convert " + key))
	return false
}

func GetEnvAsInt(key string) int {
	valStr := GetEnv(key)
	if value, err := strconv.Atoi(valStr); err == nil {
		return value
	}

	log.Fatal(errors.New("GetEnvAsInt: not convert " + key))
	return 0
}
