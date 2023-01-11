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
	// Check metric type
	switch metric.MType {
	case domain.MetricTypeGauge:
		// Convert float64 to string
		s := strconv.FormatFloat(*metric.Value, 'f', 3, 64)
		w.Write([]byte(s))
	case domain.MetricTypeCounter:
		// Convert int64 to string
		s := strconv.FormatInt(*metric.Delta, 10)
		w.Write([]byte(s))
	default:
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
		return
	}
	w.WriteHeader(http.StatusOK)
}
