package handler

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
)

func (h *Handler) mainPage(w http.ResponseWriter, r *http.Request) {
	agentID, _, _ := net.SplitHostPort(r.RemoteAddr)
	fmt.Printf("Got main page request. Method=%s Path: %s agentID: %s \n", r.Method, r.URL.Path, agentID)

	allMetrics, ok := h.s.GetAllMetrics(agentID)
	if ok {
		b := new(bytes.Buffer)
		for key, value := range allMetrics {
			fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
		}
		w.Write(b.Bytes())
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
