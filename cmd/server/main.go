package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func updateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got request. Method=%s. Path=%s\n", r.Method, r.URL.Path)
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	path := r.URL.Path
	fmt.Printf("%s\n", path)
}

func main() {
	//Configuration
	//cfg := config.New()
	//Storage
	//repo := storage.New(cfg.Server)

	// Create new chi router
	r := chi.NewRouter()

	// Using built-in middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/update/{metricName}/{metricValue}", func(r chi.Router) {
		r.Post("/", updateHandler)
	})

	//http.HandleFunc("/update/", uploadHandler)
	log.Fatal(http.ListenAndServe(":8080", r))
}
