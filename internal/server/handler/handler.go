package handler

import (
	"github.com/artems723/monik/internal/server/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	s *service.Service
}

func New(s *service.Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) InitRoutes() *chi.Mux {
	// Create new chi router
	r := chi.NewRouter()

	// Using built-in middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/update", func(r chi.Router) {
		r.Post("/{metricType}/{metricName}/{metricValue}", h.updateMetric)
		r.Post("/", h.updateMetricJSON)
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/{metricType}/{metricName}", h.getValue)
		r.Post("/", h.getValueJSON)
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/", h.mainPage)
	})
	return r
}
