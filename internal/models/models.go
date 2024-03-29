package models

import (
	"github.com/go-chi/chi/v5"
	"test-server-go/internal/config"
	"test-server-go/internal/freekassa"
	"test-server-go/internal/logger"
	"test-server-go/internal/mailer"
	"test-server-go/internal/storage"
)

type Application struct {
	Config    *config.Config
	Postgres  *storage.Postgres
	Redis     *storage.Redis
	Mailer    *mailer.Mailer
	Logger    *logger.Logger
	Router    *chi.Mux
	Freekassa *freekassa.Config
}
