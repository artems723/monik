package httpserver

import (
	"context"
	"github.com/artems723/monik/internal/server/config"
	"github.com/artems723/monik/internal/server/handler"
	"github.com/artems723/monik/internal/server/service"
	"github.com/artems723/monik/internal/server/storage"
)

// Example using http server
func ExampleServer() {
	// Create and read config
	cfg := config.Config{}
	// Create repo
	repo := storage.NewMemStorage()
	// Create service
	serv := service.New(repo, cfg)
	// Create handler
	h := handler.New(serv, cfg.Key, cfg.DatabaseDSN)
	// Create server
	srv := New()
	defer srv.Shutdown(context.Background())
	// Run server
	srv.Run(cfg.Address, h.InitRoutes())
}
