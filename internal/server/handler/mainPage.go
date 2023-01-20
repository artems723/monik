package handler

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/artems723/monik/internal/server/storage"
	"log"
	"net/http"
)

func (h *Handler) mainPage(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got main page request. Method=%s, Path: %s", r.Method, r.URL.Path)
	// Get all metrics from storage
	allMetrics, err := h.s.GetAllMetrics()
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
	// Write response
	b := new(bytes.Buffer)
	for _, value := range allMetrics.Metrics {
		fmt.Fprintf(b, "%v<br>", *value)
	}
	w.Header().Set("Content-Type", "text/html")
	_, err = w.Write(b.Bytes())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
