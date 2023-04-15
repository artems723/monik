package main

import (
	"flag"
	"log"
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
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	Key            string        `env:"KEY"`
	RateLimit      int           `env:"RATE_LIMIT"`
}

func main() {
	log.Printf("Build version: %s", buildVersion)
	log.Printf("Build date: %s", buildDate)
	log.Printf("Build commit: %s", buildCommit)
	// Create and read config
	cfg := config{}
	//Parse config from flag
	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "server address.")
	flag.DurationVar(&cfg.ReportInterval, "r", 10*time.Second, "time interval in seconds for sending metrics to server.")
	flag.DurationVar(&cfg.PollInterval, "p", 2*time.Second, "time interval in seconds for updating metrics.")
	flag.StringVar(&cfg.Key, "k", "", "key for hashing")
	flag.IntVar(&cfg.RateLimit, "l", 10, "maximum number of outgoing requests to the server")
	flag.Parse()
	// Parse config from env
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Using config: Address: %s, ReportInterval: %v, PollInterval: %v, Key: %s, RateLimit: %d", cfg.Address, cfg.ReportInterval, cfg.PollInterval, cfg.Key, cfg.RateLimit)
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

	serverAddr := "http://" + cfg.Address
	// infinite loop for sending counters to server
	reportIntervalTicker := time.NewTicker(cfg.ReportInterval)
	for range reportIntervalTicker.C {
		a.SendData(serverAddr)
	}
}
