package handler

import (
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (h *Handler) ping(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got ping request. Method=%s, Path: %s", r.Method, r.URL.Path)

	err := h.s.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
