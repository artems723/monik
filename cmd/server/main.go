package main

import (
	"github.com/artems723/monik/internal/server"
	"github.com/artems723/monik/internal/server/handler"
	"github.com/artems723/monik/internal/server/service"
	"github.com/artems723/monik/internal/server/storage"
	"github.com/caarlos0/env/v6"
	"log"
)

type config struct {
	Address string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
}

func main() {
	// Create and read config
	cfg := config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Using config: Address: %s.", cfg.Address)
	// Create storage
	repo := storage.NewMemStorage()
	// Create service
	serv := service.New(repo)
	// Create handler
	h := handler.New(serv)
	// Create server
	srv := server.New()
	// Start http server
	err = srv.Run(cfg.Address, h.InitRoutes())
	if err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
