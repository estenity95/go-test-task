package router

import (
	apimw "github.com/estenity95/go-test-task/api/middleware"
	"github.com/estenity95/go-test-task/api/resource/health"
	"github.com/estenity95/go-test-task/api/resource/subscription"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(l *zerolog.Logger, v *validator.Validate, repository subscription.Repository) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(apimw.ContentTypeJSON)

	api := subscription.New(l, v, repository)

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Get("/healthz", health.Read)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/subscriptions", func(r chi.Router) {
			r.Get("/", api.List)
			r.Post("/", api.Create)
			r.Get("/{id}", api.Read)
			r.Put("/{id}", api.Update)
			r.Delete("/{id}", api.Delete)
			r.Get("/summary", api.Summary)
		})
	})
	return r
}
