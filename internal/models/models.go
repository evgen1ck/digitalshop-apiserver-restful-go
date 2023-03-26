package models

import (
	"github.com/go-chi/chi/v5"
	"test-server-go/internal/config"
	"test-server-go/internal/database"
	"test-server-go/internal/logger"
	"test-server-go/internal/mailer"
)

type Application struct {
	Config   *config.Config
	Postgres *database.Postgres
	Mailer   *mailer.Mailer
	Logger   *logger.Logger
	Router   *chi.Mux
}
