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
	"test-server-go/internal/api_v1/handlers_v1"
	"test-server-go/internal/config"
	freekassa2 "test-server-go/internal/freekassa"
	"test-server-go/internal/logger"
	"test-server-go/internal/mailer"
	"test-server-go/internal/models"
	"test-server-go/internal/storage"
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

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	if app.Config.App.Debug {
		app.Logger.NewInfo("Prometheus API will be running in debug mode on " + prometheusServer.Addr)
		go func() {
			if err := prometheusServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				app.Logger.NewError("Error starting Prometheus API", err)
			}
		}()

		app.Logger.NewInfo("Service API v1 will be running in debug mode on " + apiV1Server.Addr)
		go func() {
			if err := apiV1Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				app.Logger.NewError("Error starting Service API v1", err)
			}
		}()
		app.Logger.NewInfo("The services is ready to listen and serve on http protocol")
	} else {
		app.Logger.NewInfo("Prometheus API will be running in tls mode on " + prometheusServer.Addr)
		go func() {
			if err := prometheusServer.ListenAndServeTLS(app.Config.Tls.CertFile, app.Config.Tls.KeyFile); err != nil && err != http.ErrServerClosed {
				app.Logger.NewError("Error starting Prometheus API", err)
			}
		}()

		app.Logger.NewInfo("Service API v1 will be running in tls mode on " + apiV1Server.Addr)
		go func() {
			if err := apiV1Server.ListenAndServeTLS(app.Config.Tls.CertFile, app.Config.Tls.KeyFile); err != nil && err != http.ErrServerClosed {
				app.Logger.NewError("Error starting Service API v1", err)
			}
		}()
		app.Logger.NewInfo("The services is ready to listen and serve on https protocol")
	}

	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		app.Logger.NewWarn("Got signal SIGINT...", nil)
	case syscall.SIGTERM:
		app.Logger.NewWarn("Got signal SIGTERM...", nil)
	}

	app.Logger.NewInfo("The services is shutting down...")

	prometheusServer.Shutdown(context.Background())
	app.Logger.NewInfo("Prometheus service is shut down")

	apiV1Server.Shutdown(context.Background())
	app.Logger.NewInfo("API v1 service is shut down")
	app.Logger.NewInfo("Done")
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

	freekassaCfg := freekassa2.NewConfig(
		uint32(cfg.Payments.Freekassa.ShopId),
		cfg.Payments.Freekassa.ApiKey,
		cfg.Payments.Freekassa.FirstSecretWord,
		cfg.Payments.Freekassa.SecondSecretWord)

	balance, err := freekassa2.Balances(freekassaCfg)
	if err != nil {
		fmt.Printf("BALANCE ERROR: %v\n", err)
	}
	if len(balance) == 0 {
		fmt.Printf("BALANCE WARN: no balance\n")
	} else {
		fmt.Printf("BALANCE:\n")
		for _, bal := range balance {
			fmt.Printf("%s %.2f\n", bal.Currency, bal.Value)
		}
	}

	url := freekassa2.NewOrderUrl(freekassaCfg, 1000.10, freekassa2.CurrencyRUB, "GTA555")
	fmt.Printf("NEW ORDER URL: %s\n", url)

	application := models.Application{
		Config:    cfg,
		Postgres:  pdb,
		Redis:     rdb,
		Mailer:    mailer.NewSmtp(*cfg),
		Logger:    zapLogger,
		Router:    chi.NewRouter(),
		Freekassa: freekassaCfg,
	}

	if application.Config.App.Debug {
		fmt.Println(application.Config)
	}

	return &application
}

func setupRouter(app models.Application) {
	r := app.Router

	r.Use(middleware.RealIP)       // Using user real ip address
	r.Use(middleware.Recoverer)    // Prevents server from crashing
	r.Use(middleware.StripSlashes) // Optimizes paths
	r.Use(middleware.Logger)       // Logging
	r.Use(middleware.Compress(5))  // Supports compression

	// prometheus routes
	setupPrometheus(app, "/prometheus/metrics", strconv.Itoa(app.Config.Prometheus.Port))

	// api version 1 routes
	rs := handlers_v1.Resolver{
		App: &app,
	}
	rs.SetupRouterApiVer1("/api/v1")
}
