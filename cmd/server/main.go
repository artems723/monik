package main

import (
	"github.com/artems723/monik/internal/server"
	"github.com/artems723/monik/internal/server/handler"
	"github.com/artems723/monik/internal/server/service"
	"github.com/artems723/monik/internal/server/storage"
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

type config struct {
	Address       string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"3s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"false"`
}

func main() {
	// Create and read config
	cfg := config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalf("error parsing config file: %v", err)
	}
	log.Printf("Using config: Address: %s", cfg.Address)
	// Create storage
	repo := storage.NewMemStorage()
	// Create service
	serv := service.New(&repo)
	// Create handler
	h := handler.New(&serv)
	// Create server
	srv := server.New()

	// Create store
	if cfg.StoreFile != "" {
		store, err := service.NewStore(cfg.StoreFile, &repo)
		if err != nil {
			log.Printf("service.NewStore, error creating Store: %v", err)
		}
		// Lead data from file to storage
		store.Init(cfg.Restore)
		//Start process of storing data to file
		go store.Run(cfg.StoreInterval)
	}

	// Start http server
	err = srv.Run(cfg.Address, h.InitRoutes())
	if err != nil {
		log.Fatalf("srv.Run, error occured while running http server: %v", err)
	}
}
