package handler

import (
	"encoding/json"
	"github.com/artems723/monik/internal/server/domain"
	"github.com/artems723/monik/internal/server/service"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) updateMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")
	log.Printf("Got update request. Method=%s, Path: %s, metricType: %s, metricName: %s, metricValue: %s\n", r.Method, r.URL.Path, metricType, metricName, metricValue)
	// Check if no metric name provided in the URL
	if metricName == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	var metric *domain.Metrics
	// Check metric type
	switch domain.MetricType(metricType) {
	case domain.MetricTypeGauge:
		// Try to convert string to float64
		val, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			log.Printf("%s. Wrong value (not float64). Got: %s\n", err, metricValue)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Create new metric
		metric = domain.NewGaugeMetric(metricName, val)
	case domain.MetricTypeCounter:
		// Convert string to int64
		val, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			log.Printf("%s. Wrong value (not int64). Got: %s\n", err, metricValue)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Create new metric
		metric = domain.NewCounterMetric(metricName, val)
	default:
		http.Error(w, domain.ErrUnknownMetricType.Error(), http.StatusNotImplemented)
		return
	}
	// Write metric to service
	err := h.s.WriteMetric(metric)
	if err != nil && err != service.ErrMTypeMismatch && err != service.ErrNoValue {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err == service.ErrMTypeMismatch || err == service.ErrNoValue {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) updateMetricJSON(w http.ResponseWriter, r *http.Request) {
	var metric *domain.Metrics
	// Read JSON and store to metric struct
	err := json.NewDecoder(r.Body).Decode(&metric)
	// Check errors
	if err != nil && err != domain.ErrUnknownMetricType {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err == domain.ErrUnknownMetricType {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}
	log.Printf("Got update JSON request. Method=%s, Path: %s, metric: %v\n", r.Method, r.URL.Path, metric)
	// Write metric to service
	err = h.s.WriteMetric(metric)
	if err != nil && err != service.ErrMTypeMismatch && err != service.ErrNoValue {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err == service.ErrMTypeMismatch || err == service.ErrNoValue {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// Encode to JSON and write to response
	err = json.NewEncoder(w).Encode(metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
