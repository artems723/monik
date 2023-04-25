// Package httpserver provides a simple http server
package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	httpServer *http.Server
}

func New() Server {
	return Server{}
}

func (s *Server) Run(serverAddr string, r *chi.Mux) error {
	s.httpServer = &http.Server{
		Addr:           serverAddr,
		MaxHeaderBytes: 1 << 20, // 1MB
		ReadTimeout:    10 * time.Second,
		// Set write timeout >30 sec. It is needed to run stress test for 30 sec using Vegeta https://github.com/tsenart/vegeta
		WriteTimeout: 31 * time.Second,
		Handler:      r,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) RunTLS(serverAddr string, certFile string, keyFile string, r *chi.Mux) error {
	s.httpServer = &http.Server{
		Addr:           serverAddr,
		MaxHeaderBytes: 1 << 20, // 1MB
		ReadTimeout:    10 * time.Second,
		// Set write timeout >30 sec. It is needed to run stress test for 30 sec using Vegeta https://github.com/tsenart/vegeta
		WriteTimeout: 31 * time.Second,
		Handler:      r,
	}
	return s.httpServer.ListenAndServeTLS(certFile, keyFile)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
