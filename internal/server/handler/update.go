package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net"
	"net/http"
	"strconv"
)

func (h *Handler) updateGaugeMetric(w http.ResponseWriter, r *http.Request) {
	agentID, _, _ := net.SplitHostPort(r.RemoteAddr)
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")
	fmt.Printf("Got gauge request. Method=%s Path: %s metricName: %s metricValue: %s\n", r.Method, r.URL.Path, metricName, metricValue)
	// Try to convert string to float64
	_, err := strconv.ParseFloat(metricValue, 32)
	if err != nil {
		fmt.Printf("%s. Wrong value (not float64). Got: %s\n", err, metricValue)
		http.Error(w, http.StatusText(400), 400)
		return
	}
	h.s.WriteMetric(agentID, metricName, metricValue)
}

func (h *Handler) updateCounterMetric(w http.ResponseWriter, r *http.Request) {
	agentID, _, _ := net.SplitHostPort(r.RemoteAddr)
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")
	fmt.Printf("Got counter request. Method=%s Path: %s metricName: %s metricValue: %s\n", r.Method, r.URL.Path, metricName, metricValue)
	// Convert string to int64
	val, err := strconv.ParseInt(metricValue, 10, 64)
	if err != nil {
		fmt.Printf("%s. Wrong value (not int64). Got: %s\n", err, metricValue)
		http.Error(w, http.StatusText(400), 400)
		return
	}
	var currentVal int64
	// Get current value from storage
	v, ok := h.s.GetMetric(agentID, metricName)
	if ok {
		// Convert string to int64
		currentVal, _ = strconv.ParseInt(v, 10, 64)
	} else {
		currentVal = int64(0)
	}
	// Sum counters
	newVal := val + currentVal
	// Write new value to storage
	h.s.WriteMetric(agentID, metricName, fmt.Sprintf("%v", newVal))
}

func (h *Handler) notImplemented(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Returning notImplemented")
	http.Error(w, http.StatusText(501), 501)
}

func (h *Handler) notFound(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Returning notFound")
	http.Error(w, http.StatusText(404), 404)
}
