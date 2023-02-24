package server

import (
	"errors"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"log"
	"net/http"
	"test-server-go/graph"
	"test-server-go/internal/models"
)

func Run(app models.Application) error {
	srv := handler.NewDefaultServer(
		graph.NewExecutableSchema(graph.Config{
			Resolvers: &graph.Resolver{App: &app},
		}),
	)
	//http.Handle("/api/v1", srv)
	http.Handle("/api/v1", errorHandler(srv))

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

func errorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic: %v", r)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
