package main

import (
	"github.com/artems723/monik/internal/client"
	"log"
	"time"
)

func newMonitor(pollInterval, reportInterval time.Duration, URL string, cl client.HTTPClient, agent client.Agent) {
	pollIntervalTicker := time.NewTicker(pollInterval)
	reportIntervalTicker := time.NewTicker(reportInterval)

	// infinite loop for polling counters and sending it to server
	for {
		select {
		case <-pollIntervalTicker.C:
			agent.UpdateMetrics()
			log.Printf("Got counters: %#v", agent)
		case <-reportIntervalTicker.C:
			agent.SendData(URL, cl)
		}
	}
}

func main() {
	endpoint := "127.0.0.1"
	port := "8080"
	URL := "http://" + endpoint + ":" + port
	cl := client.NewHTTPClient()
	agent := client.NewAgent()
	newMonitor(2*time.Second, 10*time.Second, URL, cl, agent)
}
