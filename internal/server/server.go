package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"test-server-go/internal/api_v1"
	"test-server-go/internal/api_v1/handlers"
	"test-server-go/internal/config"
	"test-server-go/internal/database"
	"test-server-go/internal/logger"
	"test-server-go/internal/mailer"
	"test-server-go/internal/models"
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
		Router:   chi.NewRouter(),
	}

	fmt.Println(app.Config)

	setupDefaultRouterSettings(app)

	srv := &http.Server{
		Addr:    app.Config.GetURL(),
		Handler: app.Router,
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

func setupDefaultRouterSettings(app models.Application) {
	r := app.Router

	r.Use(middleware.RealIP)       // Using user real ip address
	r.Use(middleware.Recoverer)    // Prevents server from crashing
	r.Use(middleware.StripSlashes) // Optimizes paths
	r.Use(middleware.Logger)       // Logging
	r.Use(middleware.Compress(5))  // Supports compression

	setupPrometheus(app)               // prometheus routes
	setupRouterApiVer1(app, "/api/v1") // api version 1 routes
}

func setupRouterApiVer1(app models.Application, pathPrefix string) {
	rs := handlers.Resolver{
		App: &app,
	}
	rs.SetupRouterApiVer1(pathPrefix)
}

func setupPrometheus(app models.Application) {
	rs := api_v1.Resolver{
		App: &app,
	}
	rs.SetupPrometheus()
}
