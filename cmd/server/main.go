package main

import (
	"github.com/artems723/monik/internal/server"
	"github.com/artems723/monik/internal/server/handler"
	"github.com/artems723/monik/internal/server/service"
	"github.com/artems723/monik/internal/server/storage"
	"github.com/caarlos0/env/v6"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
		log.Fatalf("error parsing config file: %v", err)
	}
	log.Printf("Using config: Address: %s", cfg.Address)
	// Create storage
	repo := storage.NewMemStorage()
	// Create service
	serv := service.New(repo)

	if cfg.StoreFile != "" {
		fileRepo := storage.NewFileStorage(cfg.StoreFile)
		go serv.RunFileStorage(fileRepo, cfg.Restore, cfg.StoreInterval)
	}

	// Create handler
	h := handler.New(serv)
	// Create server
	srv := server.New()

	// create channel for graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start http server
	go func() {
		err = srv.Run(cfg.Address, h.InitRoutes())
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("srv.Run, error occured while running http server: %v", err)
		}
	}()
	log.Printf("Server started")

	<-done
	// Shutdown http server
	err = srv.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("Server shutdown Failed:%+v", err)
	}
	err = serv.Shutdown()
	if err != nil {
		log.Fatalf("serv.Shutdown: %v", err)
	}
	log.Print("Server stopped properly")
}
