package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net"
	"net/http"
)

func (h *Handler) getValue(w http.ResponseWriter, r *http.Request) {
	agentID, _, _ := net.SplitHostPort(r.RemoteAddr)
	metricName := chi.URLParam(r, "metricName")
	fmt.Printf("Got get counter request. Method=%s Path: %s metricName: %s \n", r.Method, r.URL.Path, metricName)

	val, ok := h.s.GetMetric(agentID, metricName)
	if ok {
		w.Write([]byte(val))
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
