package main

import (
	"github.com/artems723/monik/internal/client"
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

func main() {

	const PollInterval = time.Second * 2
	const ReportInterval = time.Second * 10

	serverAddr := "http://localhost:8080"
	httpClient := client.NewHTTPClient()
	agent := client.NewAgent()
	newMonitor(PollInterval, ReportInterval, serverAddr, httpClient, agent)
}
