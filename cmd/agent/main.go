package main

import (
	"fmt"
	"github.com/artems723/monik/internal/client"
	"time"
)

func newMonitor(pollInterval, reportInterval int, URL string, cl client.HTTPClient, agent client.Agent) {
	pollIntervalTicker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	reportIntervalTicker := time.NewTicker(time.Duration(reportInterval) * time.Second)

	// infinite loop for polling counters and sending it to server
	for {
		select {
		case <-pollIntervalTicker.C:
			agent.UpdateMetrics()
			fmt.Print("Got counters: ")
			fmt.Println(agent)
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
	newMonitor(2, 10, URL, cl, agent)
}
