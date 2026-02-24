package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/swaggo/http-swagger"
	_ "github.com/username/project-name/docs" // Ganti dengan path modul Anda
	"github.com/username/project-name/internal/handler"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", handler.HealthCheck)
	
	// Integrasi Swagger UI
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	return r
}