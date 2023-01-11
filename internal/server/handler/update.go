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
	// Get client's IP address and use it as agentID
	agentID, _, _ := net.SplitHostPort(r.RemoteAddr)
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")
	log.Printf("Got update request. Method=%s, Path: %s, agentID: %s, metricType: %s, metricName: %s metricValue: %s\n", r.Method, r.URL.Path, agentID, metricType, metricName, metricValue)
	// Check if no metric name provided in the URL
	if metricName == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	// Check metric type
	switch domain.MetricType(metricType) {
	case domain.MetricTypeGauge:
		// Try to convert string to float64
		val, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			log.Printf("%s. Wrong value (not float64). Got: %s\n", err, metricValue)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		// Create new metric
		metric := domain.NewGaugeMetric(metricName, val)
		// Write metric to storage
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
		// Check for errors
		if err != nil && !errors.Is(err, storage.ErrNotFound) {
			log.Printf("storage.GetMetric: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if errors.Is(err, storage.ErrNotFound) {
			// Create new metric as soon as it doesn't exist in storage
			metric = domain.NewCounterMetric(metricName, int64(0))
		}
		// Add delta to current value
		*metric.Delta += val
		// Write metric with new value to storage
		err = h.s.WriteMetric(agentID, metric)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
		return
	}
}
