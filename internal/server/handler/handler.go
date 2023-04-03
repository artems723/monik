package handler

import (
	"errors"
	"github.com/artems723/monik/internal/server/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"html/template"
	"path/filepath"
)

type Handler struct {
	s           *service.Service
	key         string
	databaseDSN string
	tmpl        *template.Template
}

func New(s *service.Service, key string, databaseDSN string) *Handler {
	return &Handler{
		s:           s,
		key:         key,
		databaseDSN: databaseDSN,
		tmpl:        template.Must(template.ParseFiles(filepath.Join("templates", "mainPage.html"))),
	}
}

func (h *Handler) InitRoutes() *chi.Mux {
	// Create new chi router
	r := chi.NewRouter()

	// Using built-in middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.AllowContentEncoding("gzip"))
	r.Use(middleware.Compress(5))
	r.Use(middleware.Recoverer)

	r.Mount("/debug", middleware.Profiler())

	r.Route("/update", func(r chi.Router) {
		r.Post("/{metricType}/{metricName}/{metricValue}", h.updateMetric)
		r.Post("/", h.updateMetricJSON)
	})

	r.Route("/updates", func(r chi.Router) {
		r.Post("/", h.updateMetricsJSON)
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/{metricType}/{metricName}", h.getValue)
		r.Post("/", h.getValueJSON)
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/", h.mainPage)
	})

	r.Route("/ping", func(r chi.Router) {
		r.Get("/", h.ping)
	})
	return r
}

var ErrUnknownMetricType = errors.New("unknown metric type")
