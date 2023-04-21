package server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"test-server-go/internal/api_v1"
	"test-server-go/internal/api_v1/handlers_v1"
	"test-server-go/internal/config"
	"test-server-go/internal/logger"
	"test-server-go/internal/mailer"
	"test-server-go/internal/models"
	"test-server-go/internal/storage"
	"time"
)

func Run() {
	app := setupConfig()

	setupRouter(*app)

	prometheusServer := &http.Server{
		Addr:    "localhost:" + strconv.Itoa(app.Config.Prometheus.Port),
		Handler: app.Router,
	}
	apiV1Server := &http.Server{
		Addr:    "localhost:" + strconv.Itoa(app.Config.App.Port),
		Handler: app.Router,
	}

	done := make(chan struct{})
	go func() {
		shutdownServer(prometheusServer, app.Logger, "Prometheus API")
		shutdownServer(apiV1Server, app.Logger, "Service API v1")
		close(done)
	}()

	if app.Config.App.Debug {
		app.Logger.NewInfo("Prometheus API will be running in debug mode on " + prometheusServer.Addr)
		app.Logger.NewInfo("Service API v1 will be running in debug mode on " + apiV1Server.Addr)
		go func() {
			if err := prometheusServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				app.Logger.NewError("Error starting Prometheus API", err)
			}
		}()
		go func() {
			if err := apiV1Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				app.Logger.NewError("Error starting Service API v1", err)
			}
		}()
	} else {
		app.Logger.NewInfo("Prometheus API will be running in tls mode on" + prometheusServer.Addr)
		app.Logger.NewInfo("Service API v1 will be running in tls mode on" + apiV1Server.Addr)
		go func() {
			if err := prometheusServer.ListenAndServeTLS(app.Config.Tls.CertFile, app.Config.Tls.KeyFile); err != nil && err != http.ErrServerClosed {
				app.Logger.NewError("Error starting Prometheus API", err)
			}
		}()
		go func() {
			if err := apiV1Server.ListenAndServeTLS(app.Config.Tls.CertFile, app.Config.Tls.KeyFile); err != nil && err != http.ErrServerClosed {
				app.Logger.NewError("Error starting Service API v1", err)
			}
		}()
	}

	<-done
	app.Postgres.Pool.Close()
	app.Logger.Logger.Sync()
}

func setupConfig() *models.Application {
	ctx := context.Background()

	// Getting logger
	zapLogger, err := logger.NewZap()
	if err != nil {
		panic(err)
	}

	// Getting data config
	cfg, err := config.SetupYaml()
	if err != nil {
		zapLogger.NewError("Error creating data config", err)
	}

	// Getting PostgreSQL
	pdb, err := storage.NewPostgres(ctx, *cfg)
	if err != nil {
		zapLogger.NewError("Error connecting to the PostgreSQL database", err)
	}

	// Getting Redis
	rdb, err := storage.NewRedis(ctx, *cfg)
	if err != nil {
		zapLogger.NewError("Error connecting to the Redis database", err)
	}

	application := models.Application{
		Config:   cfg,
		Postgres: pdb,
		Redis:    rdb,
		Mailer:   mailer.NewSmtp(*cfg),
		Logger:   zapLogger,
		Router:   chi.NewRouter(),
	}

	if application.Config.App.Debug {
		fmt.Println(application.Config)
	}

	return &application
}

func shutdownServer(srv *http.Server, logger *logger.Logger, serviceName string) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	<-signalChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(ctx); err != nil {
		logger.NewError("Could not gracefully shutdown the "+serviceName, err)
	}

	logger.NewInfo(serviceName + " stopped")

}

func setupRouter(app models.Application) {
	r := app.Router

	r.Use(middleware.RealIP)       // Using user real ip address
	r.Use(middleware.Recoverer)    // Prevents server from crashing
	r.Use(middleware.StripSlashes) // Optimizes paths
	r.Use(middleware.Logger)       // Logging
	r.Use(middleware.Compress(5))  // Supports compression

	api_v1.SetupPrometheus(app) // prometheus routes
	rs := handlers_v1.Resolver{
		App: &app,
	}
	rs.SetupRouterApiVer1("/api/v1") // api version 1 routes
}
