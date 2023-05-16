package main

import (
	"flag"
	"github.com/artems723/monik/internal/server/grpcserver"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/artems723/monik/internal/server/config"
	"github.com/artems723/monik/internal/server/handler"
	"github.com/artems723/monik/internal/server/httpserver"
	"github.com/artems723/monik/internal/server/service"
	"github.com/artems723/monik/internal/server/storage"
	"github.com/caarlos0/env/v6"
	"golang.org/x/net/context"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	log.Printf("Build version: %s", buildVersion)
	log.Printf("Build date: %s", buildDate)
	log.Printf("Build commit: %s", buildCommit)
	// Create and read config
	// Config priority: default values -> json config file -> flag -> env
	cfg := config.Config{}
	// Parse config from flag
	flag.StringVar(&cfg.ConfigFile, "c", "", "json config file path")
	flag.StringVar(&cfg.Address, "a", ":8080", "server address.")
	flag.BoolVar(&cfg.Restore, "r", true, "bool value determines whether to load the initial values from the specified file when the server starts.")
	flag.DurationVar(&cfg.StoreInterval, "i", 3*time.Second, "time interval in seconds after which the current server readings are flushed to disk (value 0 makes recording synchronous).")
	path := filepath.Join(os.TempDir(), "devops-metrics-db.json")
	flag.StringVar(&cfg.StoreFile, "f", path, "string, file name where values are stored (empty value - disables writing to disk).")
	flag.StringVar(&cfg.Key, "k", "", "key for hashing")
	// Use -d "postgres://postgres:pass@postgres/postgres?sslmode=disable"
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "database data source name")
	flag.BoolVar(&cfg.EnableHTTPS, "s", false, "bool value determines whether to enable HTTPS.")
	pathCryptoKey := filepath.Join("crypto", "server.key")
	flag.StringVar(&cfg.CryptoKey, "crypto-key", pathCryptoKey, "string, crypto key path")
	pathCertFile := filepath.Join("crypto", "server.crt")
	flag.StringVar(&cfg.CertFile, "cert-file", pathCertFile, "string, cert file path")
	flag.StringVar(&cfg.TrustedSubnet, "t", "", "string, trusted subnet")
	flag.BoolVar(&cfg.GRPCEnabled, "grpc", false, "bool value determines whether to enable GRPC server.")
	flag.Parse()
	// Parse config from env
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalf("error parsing config file: %v", err)
	}
	//Parse config from json file
	if cfg.ConfigFile != "" {
		err := config.LoadJSONConfig(cfg.ConfigFile, &cfg)
		if err != nil {
			log.Fatalf("error parsing config file: %v", err)
		}
		// Parse config from flag
		flag.Parse()
		// Parse config from env
		err = env.Parse(&cfg)
		if err != nil {
			log.Fatalf("error parsing config file: %v", err)
		}
	}

	log.Printf("Using config: Address: %s, EnableHTTPS: %v, Restore: %v, StoreInterval: %v, StoreFile: %s, Key: %s, DatabaseDSN: %s, ConfigFile: %s, TrustedSubnet: %s", cfg.Address, cfg.EnableHTTPS, cfg.Restore, cfg.StoreInterval, cfg.StoreFile, cfg.Key, cfg.DatabaseDSN, cfg.ConfigFile, cfg.TrustedSubnet)

	// Create storage
	var repo service.Repository
	if cfg.DatabaseDSN != "" {
		repo, err = storage.NewPostgresStorage(cfg.DatabaseDSN)
		if err != nil {
			log.Fatalf("storage.NewPostgresStorage, error occured while connecting to postgres: %v", err)
		}
		// Disable file storage feature
		cfg.StoreFile = ""
	} else {
		repo = storage.NewMemStorage()
	}
	// Create service
	serv := service.New(repo, cfg)

	// Run file storage
	if cfg.StoreFile != "" {
		fileRepo := storage.NewFileStorage(cfg.StoreFile)
		go serv.RunFileStorage(fileRepo)
	}

	// Create handler
	h := handler.New(serv, cfg.Key, cfg.DatabaseDSN)
	// Create server
	srv := httpserver.New()

	if cfg.GRPCEnabled {
		grpcsrv := grpcserver.New(serv, cfg)
		go grpcsrv.Start()
	}

	// Create channel for graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start http server
	go srv.Start(cfg, h.InitRoutes(cfg))
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
