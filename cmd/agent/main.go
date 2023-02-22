package main

import (
	"flag"
	"github.com/artems723/monik/internal/client/agent"
	"github.com/artems723/monik/internal/client/httpClient"
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

type config struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	Key            string        `env:"KEY"`
	RateLimit      int           `env:"RATE_LIMIT"`
}

func main() {
	// Create and read config
	cfg := config{}
	//Parse config from flag
	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "server address.")
	flag.DurationVar(&cfg.ReportInterval, "r", 10*time.Second, "time interval in seconds after which a reports metrics to server.")
	flag.DurationVar(&cfg.PollInterval, "p", 2*time.Second, "time interval in seconds after which a updates metrics.")
	flag.StringVar(&cfg.Key, "k", "", "key for hashing")
	flag.IntVar(&cfg.RateLimit, "l", 10, "maximum number of outgoing requests to the server")
	flag.Parse()
	// Parse config from env
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Using config: Address: %s, ReportInterval: %v, PollInterval: %v, Key: %s, RateLimit: %d", cfg.Address, cfg.ReportInterval, cfg.PollInterval, cfg.Key, cfg.RateLimit)

	serverAddr := "http://" + cfg.Address

	cl := httpClient.NewHTTPClient()
	a := agent.NewAgent(cfg.Key, cl)

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

	if cfg.RateLimit <= 0 {
		log.Fatal("RateLimit must be greater than 0")
	}

	// use workerpool pattern to limit maximum number of outgoing connections to server
	//jobCh := make(chan struct{})
	//for i := 0; i < cfg.RateLimit; i++ {
	//	go func() {
	//		for range jobCh {
	//			a.SendData(serverAddr)
	//		}
	//	}()
	//}
	// infinite loop for sending counters to server
	reportIntervalTicker := time.NewTicker(cfg.ReportInterval)
	for range reportIntervalTicker.C {
		a.SendData(serverAddr)
		//jobCh <- struct{}{}
	}
}
