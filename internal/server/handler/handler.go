package handler

import (
	"github.com/artems723/monik/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

type Handler struct {
	s storage.Repository
}

func New(s storage.Repository) Handler {
	return Handler{s: s}
}

func (h *Handler) InitRoutes() *chi.Mux {
	// Create new chi router
	r := chi.NewRouter()

	// Using built-in middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Route /update path
	r.Route("/update", func(r chi.Router) {
		r.Post("/gauge/{metricName}/{metricValue}", h.updateGaugeMetric)
		r.Post("/counter/{metricName}/{metricValue}", h.updateCounterMetric)
		r.Post("/gauge/", http.NotFound)
		r.Post("/counter/", http.NotFound)
		r.Post("/*", h.notImplemented)
	})
	return r
}
