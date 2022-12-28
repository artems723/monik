package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handler) updateMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")
	fmt.Printf("Got request. Method=%s Path: %s metricType: %s metricName: %s metricValue: %s\n", r.Method, r.URL.Path, metricType, metricName, metricValue)
	h.s.Write(metricType, metricName, metricValue)
}
