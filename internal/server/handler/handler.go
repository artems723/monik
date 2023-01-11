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
		r.Post("/{metricType}/{metricName}/{metricValue}", h.updateMetric)
		r.Post("/gauge/", http.NotFound)
		r.Post("/counter/", http.NotFound)
		r.Post("/*", h.notImplemented)
	})

	// Route /value path
	r.Route("/value", func(r chi.Router) {
		r.Get("/{metricType}/{metricName}", h.getValue)
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/", h.mainPage)
	})
	return r
}
