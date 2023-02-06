package server

import (
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"net/http"
	"test-server-go/graph"
	"test-server-go/internal/config"
	"test-server-go/internal/mailer"
	"test-server-go/internal/postgres"
)

type Application struct {
	Config   *config.Config
	Postgres *postgres.Postgres
	Mailer   *mailer.Mailer
	//logger       *logger.Logger
	//sessionStore *sessions.CookieStore
}

func (app *Application) ServerRun() error {
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	http.Handle("/query", srv)

	if app.Config.App.DebugMode == true {
		fmt.Println("server is running in debug mode")
		http.Handle("/", playground.Handler("GraphQL playground", "/query"))
		err := http.ListenAndServe(app.Config.GetURL(), nil)
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	} else {
		fmt.Println("server is running in tls mode")
		err := http.ListenAndServeTLS(app.Config.GetURL(), app.Config.Tls.CertFile, app.Config.Tls.KeyFile, nil)
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}

	return nil
}
