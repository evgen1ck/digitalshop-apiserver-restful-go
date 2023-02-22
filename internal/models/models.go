package models

import (
	"test-server-go/internal/config"
	"test-server-go/internal/logger"
	"test-server-go/internal/mailer"
	"test-server-go/internal/postgres"
)

type Application struct {
	Config   *config.Config
	Postgres *postgres.Postgres
	Mailer   *mailer.Mailer
	Logrus   *logger.Logger
}

type ExistsNicknameEmail struct {
	NicknameExists bool
	EmailExists    bool
}
