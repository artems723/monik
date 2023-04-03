package main

import (
	"flag"
	"github.com/artems723/monik/internal/server"
	"github.com/artems723/monik/internal/server/handler"
	"github.com/artems723/monik/internal/server/service"
	"github.com/artems723/monik/internal/server/storage"
	"github.com/caarlos0/env/v6"
	"golang.org/x/net/context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Create and read config
	cfg := server.Config{}
	//Parse config from flag
	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "server address.")
	flag.BoolVar(&cfg.Restore, "r", true, "bool value determines whether to load the initial values from the specified file when the server starts.")
	flag.DurationVar(&cfg.StoreInterval, "i", 3*time.Second, "time interval in seconds after which the current server readings are flushed to disk (value 0 makes recording synchronous).")
	flag.StringVar(&cfg.StoreFile, "f", "/tmp/devops-metrics-db.json", "string, file name where values are stored (empty value - disables writing to disk).")
	flag.Parse()
	// Parse config from env
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalf("error parsing config file: %v", err)
	}
	log.Printf("Using config: Address: %s, Restore: %v, StoreInterval: %v, StoreFile: %s", cfg.Address, cfg.Restore, cfg.StoreInterval, cfg.StoreFile)
	// Create storage
	repo := storage.NewMemStorage()
	// Create service
	serv := service.New(repo, cfg)

	if cfg.StoreFile != "" {
		fileRepo := storage.NewFileStorage(cfg.StoreFile)
		go serv.RunFileStorage(fileRepo)
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
