package server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test-server-go/internal/api_v1"
	"test-server-go/internal/api_v1/handlers"
	"test-server-go/internal/config"
	"test-server-go/internal/database"
	"test-server-go/internal/logger"
	"test-server-go/internal/mailer"
	"test-server-go/internal/models"
	"time"
)

func Run() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config build error %s", err.Error())
	}

	zapLogger, err := logger.New("zip")
	if err != nil {
		panic(err)
	}
	defer zapLogger.Sync()

	pdb, err := database.NewPostgres(context.Background(), cfg.GetPostgresDSN())
	if err != nil {
		zapLogger.NewError("Error connecting to the database database", err)
	}

	app := models.Application{
		Config:   cfg,
		Postgres: pdb,
		Mailer:   mailer.NewSmtp(*cfg),
		Logger:   zapLogger,
		Router:   chi.NewRouter(),
	}

	fmt.Println(app.Config)

	setupRouterSettings(app)

	srv := &http.Server{
		Addr:    app.Config.GetLocalUrlApp(),
		Handler: app.Router,
	}

	go shutdownServer(srv, zapLogger)

	if app.Config.App.Debug {
		app.Logger.NewInfo("Server is running in debug mode")
		_ = srv.ListenAndServe()
	} else {
		app.Logger.NewInfo("Server is running in tls mode")
		_ = srv.ListenAndServeTLS(app.Config.Tls.CertFile, app.Config.Tls.KeyFile)
	}
}

func shutdownServer(srv *http.Server, logger *logger.Logger) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	<-signalChan

	logger.NewInfo("Server is shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(ctx); err != nil {
		logger.NewError("Could not gracefully shutdown the server", err)
	}

	logger.NewInfo("Server stopped")
}

func setupRouterSettings(app models.Application) {
	r := app.Router

	r.Use(middleware.RealIP)       // Using user real ip address
	r.Use(middleware.Recoverer)    // Prevents server from crashing
	r.Use(middleware.StripSlashes) // Optimizes paths
	r.Use(middleware.Logger)       // Logging
	r.Use(middleware.Compress(5))  // Supports compression

	api_v1.SetupPrometheus(app) // prometheus routes
	rs := handlers.Resolver{
		App: &app,
	}
	rs.SetupRouterApiVer1("/api/v1") // api version 1 routes
}
