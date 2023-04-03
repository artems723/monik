package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/artems723/monik/internal/server/domain"
	"github.com/artems723/monik/internal/server/service"
	"github.com/artems723/monik/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) getValue(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	log.Printf("Got get value request. Method=%s, Path: %s, metricType: %s, metricName: %s\n", r.Method, r.URL.Path, metricType, metricName)
	// Get metric from service
	metric, err := h.s.GetMetric(r.Context(), domain.NewMetric(metricName, metricType))
	// Check for errors
	if err != nil && !errors.Is(err, storage.ErrNotFound) && !errors.Is(err, service.ErrMTypeMismatch) {
		log.Printf("storage.GetMetric: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if errors.Is(err, storage.ErrNotFound) || errors.Is(err, service.ErrMTypeMismatch) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	var str string
	// Check metric type
	switch metric.MType {
	case domain.MetricTypeGauge:
		// Convert float64 to string
		str = strconv.FormatFloat(*metric.Value, 'f', 3, 64)
	case domain.MetricTypeCounter:
		// Convert int64 to string
		str = strconv.FormatInt(*metric.Delta, 10)
	case domain.MetricTypeUnknown:
		http.Error(w, ErrUnknownMetricType.Error(), http.StatusNotImplemented)
		return
	}
	_, err = w.Write([]byte(str))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getValueJSON(w http.ResponseWriter, r *http.Request) {
	var metric *domain.Metric
	// Read JSON and store to metric struct
	err := json.NewDecoder(r.Body).Decode(&metric)
	// Check errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Got get value JSON request. Method=%s, Path: %s, metric: %v\n", r.Method, r.URL.Path, metric)
	if metric.MType == domain.MetricTypeUnknown {
		http.Error(w, ErrUnknownMetricType.Error(), http.StatusNotImplemented)
		return
	}
	// Get metric from service
	res, err := h.s.GetMetric(r.Context(), metric)
	// Check for errors
	if err != nil && !errors.Is(err, storage.ErrNotFound) && !errors.Is(err, service.ErrMTypeMismatch) {
		log.Printf("storage.GetMetric: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if errors.Is(err, storage.ErrNotFound) || errors.Is(err, service.ErrMTypeMismatch) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Add hash to metric if key was provided
	if h.key != "" {
		res.AddHash(h.key)
	}

	w.Header().Set("Content-Type", "application/json")
	// Encode to JSON and write to response
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
