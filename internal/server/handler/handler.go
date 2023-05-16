// Package handler contains all handlers for server
package handler

import (
	"errors"
	"github.com/artems723/monik/internal/server/config"
	"html/template"
	"log"
	"net"
	"net/http"
	"path/filepath"

	"github.com/artems723/monik/internal/server/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	databaseDSN string
	key         string
	s           *service.Service
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

func (h *Handler) InitRoutes(cfg config.Config) *chi.Mux {
	// Create new chi router
	r := chi.NewRouter()

	// Check if real ip is from trusted subnet
	r.Use(CheckTrustedSubnet(cfg.TrustedSubnet))
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

func CheckTrustedSubnet(trustedSubnet string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if trustedSubnet == "" {
				next.ServeHTTP(w, r)
				return
			}
			realIP := r.Header.Get("X-Real-IP")
			if realIP == "" {
				http.Error(w, "headers not contain X-Real-IP", http.StatusBadRequest)
				return
			}
			ip := net.ParseIP(realIP)
			if ip == nil {
				http.Error(w, "cannot parse X-Real-IP", http.StatusBadRequest)
				return
			}
			_, subnet, err := net.ParseCIDR(trustedSubnet)
			if err != nil {
				log.Printf("cannot parse trusted subnet: %s", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			if !subnet.Contains(ip) {
				http.Error(w, "ip not from trusted subnet", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
