package server

import (
	"errors"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"net/http"
	"test-server-go/graph"
	"test-server-go/internal/models"
)

func Run(app models.Application) error {
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	http.Handle("/api/v1", srv)

	if app.Config.App.DebugMode {
		app.Logrus.NewInfo("Server is running in debug mode")
		http.Handle("/", playground.Handler("GraphQL playground", "/api/v1"))
		err := http.ListenAndServe(app.Config.GetURL(), nil)
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	} else {
		app.Logrus.NewInfo("Server is running in tls mode")
		err := http.ListenAndServeTLS(app.Config.GetURL(), app.Config.Tls.CertFile, app.Config.Tls.KeyFile, nil)
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}

	return nil
}
