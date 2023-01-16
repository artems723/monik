package main

import (
	"github.com/artems723/monik/internal/client"
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

func newMonitor(pollInterval, reportInterval time.Duration, serverAddr string, httpClient client.HTTPClient, agent client.Agent) {
	pollIntervalTicker := time.NewTicker(pollInterval)
	reportIntervalTicker := time.NewTicker(reportInterval)

	// infinite loop for polling counters and sending it to server
	for {
		select {
		case <-pollIntervalTicker.C:
			agent.UpdateMetrics()
		case <-reportIntervalTicker.C:
			agent.SendData(serverAddr, httpClient)
		}
	}
}

type config struct {
	Address        string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
}

func main() {
	// Create and read config
	cfg := config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Using config: Address: %s, ReportInterval: %v, PollInterval: %v.", cfg.Address, cfg.ReportInterval, cfg.PollInterval)

	serverAddr := "http://" + cfg.Address

	httpClient := client.NewHTTPClient()
	agent := client.NewAgent()
	newMonitor(cfg.PollInterval, cfg.ReportInterval, serverAddr, httpClient, agent)
}
