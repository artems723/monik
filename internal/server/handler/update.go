package handler

import (
	"errors"
	"github.com/artems723/monik/internal/server/domain"
	"github.com/artems723/monik/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"net"
	"net/http"
	"strconv"
)

func (h *Handler) updateMetric(w http.ResponseWriter, r *http.Request) {
	agentID, _, _ := net.SplitHostPort(r.RemoteAddr)
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")
	log.Printf("Got update request. Method=%s, Path: %s, agentID: %s, metricType: %s, metricName: %s metricValue: %s\n", r.Method, r.URL.Path, agentID, metricType, metricName, metricValue)

	switch domain.MetricType(metricType) {
	case domain.MetricTypeGauge:
		// Try to convert string to float64
		val, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			log.Printf("%s. Wrong value (not float64). Got: %s\n", err, metricValue)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		metric := domain.NewGaugeMetric(metricName, val)
		err = h.s.WriteMetric(agentID, metric)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	case domain.MetricTypeCounter:
		// Convert string to int64
		val, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			log.Printf("%s. Wrong value (not int64). Got: %s\n", err, metricValue)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		// Get current metric from storage
		metric, err := h.s.GetMetric(agentID, metricName)
		if err != nil && !errors.Is(err, storage.ErrNotFound) {
			log.Printf("storage.GetMetric: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if errors.Is(err, storage.ErrNotFound) {
			metric = domain.NewCounterMetric(metricName, int64(0))
		}
		// Sum counters
		*metric.Delta += val
		// Write new value to storage
		err = h.s.WriteMetric(agentID, metric)
		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
		return
	}
}

func (h *Handler) notImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
