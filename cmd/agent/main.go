package main

import (
	"fmt"
	"github.com/artems723/monik/internal/client"
	"time"
)

func newMonitor(pollInterval, reportInterval int, endpoint string, port int, cl client.Client, agent client.Agent) {
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
			agent.SendData(endpoint, port, cl)
		}
	}
}

func main() {
	endpoint := "127.0.0.1"
	port := 8080
	cl := client.New()
	agent := client.NewAgent()
	newMonitor(2, 2, endpoint, port, cl, agent)
}
