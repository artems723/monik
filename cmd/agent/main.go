package main

import (
	"flag"
	"log"
	"path/filepath"
	"time"

	"github.com/artems723/monik/internal/client/agent"
	"github.com/artems723/monik/internal/client/httpclient"
	"github.com/caarlos0/env/v6"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

type config struct {
	Address        string        `env:"ADDRESS"`
	ConfigFile     string        `env:"CONFIG"`
	CryptoKey      string        `env:"CRYPTO_KEY"`
	EnableHTTPS    bool          `env:"ENABLE_HTTPS"`
	Key            string        `env:"KEY"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	RateLimit      int           `env:"RATE_LIMIT"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
}

func main() {
	log.Printf("Build version: %s", buildVersion)
	log.Printf("Build date: %s", buildDate)
	log.Printf("Build commit: %s", buildCommit)
	// Create and read config
	cfg := config{}
	//Parse config from flag
	flag.StringVar(&cfg.ConfigFile, "c", "", "json config file path")
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "server address.")
	flag.DurationVar(&cfg.ReportInterval, "r", 10*time.Second, "time interval in seconds for sending metrics to server.")
	flag.DurationVar(&cfg.PollInterval, "p", 2*time.Second, "time interval in seconds for updating metrics.")
	flag.StringVar(&cfg.Key, "k", "", "key for hashing")
	flag.IntVar(&cfg.RateLimit, "l", 10, "maximum number of outgoing requests to the server")
	pathCryptoKey := filepath.Join("crypto", "server.crt")
	flag.StringVar(&cfg.CryptoKey, "crypto-key", pathCryptoKey, "string, crypto key path")
	flag.BoolVar(&cfg.EnableHTTPS, "s", false, "bool value determines whether to enable HTTPS.")
	flag.Parse()
	// Parse config from env
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Using config: Address: %s, EnableHTTPS: %v, ReportInterval: %v, PollInterval: %v, Key: %s, RateLimit: %d", cfg.Address, cfg.EnableHTTPS, cfg.ReportInterval, cfg.PollInterval, cfg.Key, cfg.RateLimit)
	if cfg.RateLimit <= 0 {
		log.Fatal("RateLimit must be greater than 0")
	}
	// create agent with httpClient
	cl := httpclient.New(cfg.RateLimit)
	a := agent.New(cfg.Key, cl)

	// infinite loop for polling counters
	pollIntervalTicker1 := time.NewTicker(cfg.PollInterval)
	go func() {
		for range pollIntervalTicker1.C {
			a.UpdateMetrics()
		}
	}()
	// infinite loop for polling additional counters
	pollIntervalTicker2 := time.NewTicker(cfg.PollInterval)
	go func() {
		for range pollIntervalTicker2.C {
			a.UpdateAdditionalMetrics()
		}
	}()

	var serverAddr string
	switch cfg.EnableHTTPS {
	case false:
		serverAddr = "http://" + cfg.Address
	case true:
		serverAddr = "https://" + cfg.Address
		cl.SetRootCertificate(cfg.CryptoKey)
	}
	// infinite loop for sending counters to server
	reportIntervalTicker := time.NewTicker(cfg.ReportInterval)
	for range reportIntervalTicker.C {
		a.SendData(serverAddr)
	}
}
