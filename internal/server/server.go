package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"test-server-go/internal/api"
	"test-server-go/internal/config"
	"test-server-go/internal/database"
	"test-server-go/internal/logger"
	"test-server-go/internal/mailer"
	"test-server-go/internal/models"
	"time"
)

func Run() {
	logrus := logger.New()

	cfg, err := config.New(logrus)
	if err != nil {
		logrus.NewErrorWithExit("Config build error", err)
	}

	pdb, err := database.NewPostgres(context.Background(), cfg.GetPostgresDSN())
	if err != nil {
		logrus.NewErrorWithExit("Error connecting to the database database", err)
	}
	defer pdb.Close()

	app := models.Application{
		Config:   cfg,
		Postgres: pdb,
		Mailer:   mailer.NewSmtp(*cfg),
		Logrus:   logrus,
	}

	fmt.Println(app.Config)

	resolver22 := api.Resolver22{
		App: &app,
	}
	router3 := resolver22.NewRoutes()

	srv := &http.Server{
		Addr:         app.Config.GetURL(),
		Handler:      router3,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	if app.Config.App.DebugMode {
		app.Logrus.NewInfo("Server is running in debug mode")
		err := srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			app.Logrus.NewErrorWithExit("Server crashed", err)
		}
	} else {
		app.Logrus.NewInfo("Server is running in tls mode")
		err := srv.ListenAndServeTLS(app.Config.Tls.CertFile, app.Config.Tls.KeyFile)
		if !errors.Is(err, http.ErrServerClosed) {
			app.Logrus.NewErrorWithExit("Server crashed", err)
		}
	}
}
