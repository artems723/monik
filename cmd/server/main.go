package main

import (
	"github.com/artems723/monik/internal/server"
	"github.com/artems723/monik/internal/server/handler"
	"github.com/artems723/monik/internal/server/storage"
	"log"
)

func main() {
	// Create storage
	repo := storage.NewMemStorage()
	// Create handler
	h := handler.New(repo)
	// Create server
	srv := server.New()
	// Start http server
	err := srv.Run("8080", h.InitRoutes())
	if err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
