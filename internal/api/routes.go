package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"net/http"
	"test-server-go/internal/models"
	"time"
)

type Resolver2 struct {
	App *models.Application
}

func (rd *Resolver2) NewRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.StripSlashes)                    // Оптимизация путей
	r.Use(middleware.Compress(5, "application/json")) // Поддержка сжатия
	r.Use(middleware.Logger)                          // Логгирование
	r.Use(middleware.Recoverer)                       // Обработка ошибок 500
	r.Use(httprate.Limit(                             // Количество запросов для одного ip (10 в 10 секунд)
		10,
		10*time.Second,
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			respondWithTooManyRequests(w)
		}),
	))
	r.NotFound(func(w http.ResponseWriter, r *http.Request) { // Обработка ошибок 404
		respondWithResourceNotFound(w)
	})
	r.Use(cors.New(cors.Options{ // Разрешаем CORS
		AllowedOrigins:   []string{"*"},                                                       // Разрешаем запросы со всех источников
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},                 // Разрешаем все методы запросов
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}, // Разрешаем указанные заголовки
		ExposedHeaders:   []string{"Link"},                                                    // Разрешаем доступ к указанным заголовкам из JavaScript кода
		AllowCredentials: true,                                                                // Разрешаем куки в CORS запросах
		MaxAge:           300,                                                                 // Устанавливаем максимальный срок действия CORS политики
	}).Handler)

	r.Mount("/api/v1", rd.Routes(r))

	return r
}

func (rd *Resolver2) Routes(r *chi.Mux) http.Handler {
	r.Get("/getAllAlbums3", getAllAlbums)
	r.Get("/getAllAlbums", getAllAlbums)
	r.Get("/authSignup", rd.authSignup)

	return r
}
