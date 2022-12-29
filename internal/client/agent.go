package client

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"net/http"
	"runtime"
)

// metric types
type (
	gauge   float64
	counter int64

	Agent struct {
		id           string
		gaugeMetrics map[string]gauge
		pollCount    counter
	}
)

func NewAgent() Agent {
	id := uuid.New().String()
	metricsMap := make(map[string]gauge)
	return Agent{id: id, gaugeMetrics: metricsMap, pollCount: 0}
}

func (agent *Agent) UpdateMetrics() {
	var rtm runtime.MemStats

	//Read memory stats
	runtime.ReadMemStats(&rtm)

	//Update metrics
	agent.gaugeMetrics["Alloc"] = gauge(rtm.Alloc)
	agent.gaugeMetrics["BuckHashSys"] = gauge(rtm.BuckHashSys)
	agent.gaugeMetrics["Frees"] = gauge(rtm.Frees)
	agent.gaugeMetrics["GCCPUFraction"] = gauge(rtm.GCCPUFraction)
	agent.gaugeMetrics["GCSys"] = gauge(rtm.GCSys)
	agent.gaugeMetrics["HeapAlloc"] = gauge(rtm.HeapAlloc)
	agent.gaugeMetrics["HeapIdle"] = gauge(rtm.HeapIdle)
	agent.gaugeMetrics["HeapInuse"] = gauge(rtm.HeapInuse)
	agent.gaugeMetrics["HeapObjects"] = gauge(rtm.HeapObjects)
	agent.gaugeMetrics["HeapReleased"] = gauge(rtm.HeapReleased)
	agent.gaugeMetrics["HeapSys"] = gauge(rtm.HeapSys)
	agent.gaugeMetrics["LastGC"] = gauge(rtm.LastGC)
	agent.gaugeMetrics["Lookups"] = gauge(rtm.Lookups)
	agent.gaugeMetrics["MCacheInuse"] = gauge(rtm.MCacheInuse)
	agent.gaugeMetrics["MCacheSys"] = gauge(rtm.MCacheSys)
	agent.gaugeMetrics["MSpanInuse"] = gauge(rtm.MSpanInuse)
	agent.gaugeMetrics["MSpanSys"] = gauge(rtm.MSpanSys)
	agent.gaugeMetrics["Mallocs"] = gauge(rtm.Mallocs)
	agent.gaugeMetrics["NextGC"] = gauge(rtm.NextGC)
	agent.gaugeMetrics["NumForcedGC"] = gauge(rtm.NumForcedGC)
	agent.gaugeMetrics["NumGC"] = gauge(rtm.NumGC)
	agent.gaugeMetrics["OtherSys"] = gauge(rtm.OtherSys)
	agent.gaugeMetrics["PauseTotalNs"] = gauge(rtm.PauseTotalNs)
	agent.gaugeMetrics["StackInuse"] = gauge(rtm.StackInuse)
	agent.gaugeMetrics["StackSys"] = gauge(rtm.StackSys)
	agent.gaugeMetrics["Sys"] = gauge(rtm.Sys)
	agent.gaugeMetrics["TotalAlloc"] = gauge(rtm.TotalAlloc)
	agent.pollCount++
	agent.gaugeMetrics["RandomValue"] = gauge(rand.Float64())
}

// send metrics data to http server
func (agent *Agent) SendData(endpoint string, port int, client Client) {

	// send gauges
	for key, val := range agent.gaugeMetrics {
		urlString := fmt.Sprintf("http://%s:%d/update/gauge/%s/%f", endpoint, port, key, val)

		fmt.Printf("Sending data to %s\n", urlString)
		resp, err := client.client.R().SetHeader("Content-Type", "text/plain").Post(urlString)
		if err != nil {
			fmt.Printf("Error sending request: %s\n", err)
			return
		}
		fmt.Printf("The status code we got is: %d %s\n", resp.StatusCode, http.StatusText(resp.StatusCode()))
	}

	// send counter
	urlString := fmt.Sprintf("http://%s:%d/update/counter/PollCount/%d", endpoint, port, agent.pollCount)
	resp, err := client.client.R().SetHeader("Content-Type", "text/plain").Post(urlString)
	if err != nil {
		fmt.Printf("Error sending request: %s\n", err)
		return
	}
	fmt.Printf("The status code we got is: %d %s\n", resp.StatusCode, http.StatusText(resp.StatusCode()))

	// reset the counter
	agent.pollCount = 0
}
