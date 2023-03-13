package models

import (
	"test-server-go/internal/config"
	"test-server-go/internal/database"
	"test-server-go/internal/logger"
	"test-server-go/internal/mailer"
)

type Application struct {
	Config   *config.Config
	Postgres *database.Postgres
	Mailer   *mailer.Mailer
	Logrus   *logger.Logger
}
