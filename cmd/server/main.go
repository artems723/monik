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
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
}

func main() {
	// Create and read config
	cfg := config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Using config: Address: %s", cfg.Address)
	// Create storage
	repo := storage.NewMemStorage()
	// Create service
	serv := service.New(repo)
	// Create handler
	h := handler.New(serv)
	// Create server
	srv := server.New()
	// Create store
	//store, err := server.NewStore(cfg.StoreFile)
	//if err != nil {
	//	log.Fatalf("error occured while creating store: %s", err.Error())
	//}
	//metrics, err := store.ReadMetrics()
	//if err != nil {
	//	log.Fatalf("error occured while reading metrics from file: %s", err.Error())
	//}
	//err = serv.WriteMetrics(metrics)
	//if err != nil {
	//	log.Fatalf("error occured while writing metrics to storage: %s", err.Error())
	//}
	// Start http server
	err = srv.Run(cfg.Address, h.InitRoutes())
	if err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
	// Start process of storing data to file
	//err = store.Run(cfg.StoreInterval)
	if err != nil {
		log.Fatalf("error occured while running store: %s", err.Error())
	}
}
