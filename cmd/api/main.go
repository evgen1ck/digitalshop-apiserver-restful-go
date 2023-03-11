package main

import (
	"context"
	"fmt"
	"test-server-go/internal/config"
	"test-server-go/internal/database"
	"test-server-go/internal/logger"
	"test-server-go/internal/mailer"
	"test-server-go/internal/models"
	"test-server-go/internal/server"
)

func main() {
	logrus := logger.New()

	cfg, err := config.New(logrus)
	if err != nil {
		logrus.NewError("Config build error", err)
	}

	pdb, err := database.NewPostgres(context.Background(), cfg.GetPostgresDSN())
	if err != nil {
		logrus.NewError("Error connecting to the database database", err)
	}
	defer pdb.Close()

	application := models.Application{
		Config:   cfg,
		Postgres: pdb,
		Mailer:   mailer.NewSmtp(*cfg),
		Logrus:   logrus,
	}

	fmt.Println(application.Config)

	err = server.Run(application)
	if err != nil {
		logrus.NewError("Server startup error", err)
	}

	logrus.NewInfo("Server is stopped")
}
