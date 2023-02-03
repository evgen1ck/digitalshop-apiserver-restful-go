package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	graph2 "test-server-go/graph"
	"test-server-go/internal/config"
	"test-server-go/internal/mailer"
	"test-server-go/internal/postgres"
	"time"
)

type Application struct {
	Config   *config.Config
	Postgres *postgres.Postgres
	Mailer   *mailer.Mailer
	//logger       *logger.Logger
	//sessionStore *sessions.CookieStore
}

func (app *Application) ServerRun() error {

	apiGraphql := handler.NewDefaultServer(graph2.NewExecutableSchema(graph2.Config{Resolvers: &graph2.Resolver{}}))
	http.Handle("/query", apiGraphql)

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         app.Config.GetURL(),
		Handler:      apiGraphql,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		TLSConfig:    tlsConfig,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	if app.Config.App.DebugMode == true {
		fmt.Println("server is running in debug mode")
		http.Handle("/", playground.Handler("GraphQL playground", "/query"))
		err := srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	} else {
		fmt.Println("server is running in tls mode")
		err := srv.ListenAndServeTLS(app.Config.Tls.CertFile, app.Config.Tls.KeyFile)
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}

	return <-shutdownError
}
