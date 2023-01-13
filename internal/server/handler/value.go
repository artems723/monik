package handler

import (
	"encoding/json"
	"errors"
	"github.com/artems723/monik/internal/server/domain"
	"github.com/artems723/monik/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"net"
	"net/http"
	"strconv"
)

func (h *Handler) getValue(w http.ResponseWriter, r *http.Request) {
	// Get client's IP address and use it as agentID
	agentID, _, _ := net.SplitHostPort(r.RemoteAddr)
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	log.Printf("Got get value request. Method=%s, Path: %s, agentID: %s, metricType: %s, metricName: %s\n", r.Method, r.URL.Path, agentID, metricType, metricName)
	// Get metric from storage
	metric, err := h.s.GetMetric(agentID, metricName)
	// Check for errors
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		log.Printf("storage.GetMetric: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if errors.Is(err, storage.ErrNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
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
	default:
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
		return
	}
	w.Write([]byte(str))
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getValueJSON(w http.ResponseWriter, r *http.Request) {
	// Get client's IP address and use it as agentID
	agentID, _, _ := net.SplitHostPort(r.RemoteAddr)

	var metric domain.Metrics
	// Read JSON and store to metric struct
	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Got get value JSON request. Method=%s, Path: %s, agentID: %s, metricType: %s, metricName: %s\n", r.Method, r.URL.Path, agentID, metric.MType, metric.ID)

	// TODO: create service
	// Get metric from storage
	res, err := h.s.GetMetric(agentID, metric.ID)
	// Check for errors
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		log.Printf("storage.GetMetric: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if errors.Is(err, storage.ErrNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	// Encode to JSON
	resJson, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resJson)
}
