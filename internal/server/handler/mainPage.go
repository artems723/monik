package handler

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/artems723/monik/internal/server/storage"
	"log"
	"net"
	"net/http"
)

func (h *Handler) mainPage(w http.ResponseWriter, r *http.Request) {
	// Get client's IP address and use it as agentID
	agentID, _, _ := net.SplitHostPort(r.RemoteAddr)
	log.Printf("Got main page request. Method=%s, Path: %s, agentID: %s\n", r.Method, r.URL.Path, agentID)
	// Get all metrics from storage
	allMetrics, err := h.s.GetAllMetrics(agentID)
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
	for key, value := range allMetrics {
		fmt.Fprintf(b, "%s=\"%v\"\n", key, value.String())
	}
	w.Write(b.Bytes())
	w.WriteHeader(http.StatusOK)
}
