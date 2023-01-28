package handler

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"net/http"
	"time"
)

func (h *Handler) ping(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got ping request. Method=%s, Path: %s", r.Method, r.URL.Path)

	db, err := sql.Open("pgx", h.databaseDSN)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
