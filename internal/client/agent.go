package client

import (
	"fmt"
	"math/rand"
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
	metricsMap := make(map[string]gauge)
	return Agent{gaugeMetrics: metricsMap, pollCount: 0}
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
func (agent *Agent) SendData(URL string, client HttpClient) error {

	// send gauges
	for key, val := range agent.gaugeMetrics {
		urlString := fmt.Sprintf("%s/update/gauge/%s/%f", URL, key, val)

		fmt.Printf("Sending data to %s\n", urlString)
		_, err := client.client.R().SetHeader("Content-Type", "text/plain").Post(urlString)
		if err != nil {
			fmt.Printf("Error sending request: %s\n", err)
			return err
		}
	}

	// send counter
	urlString := fmt.Sprintf("%s/update/counter/PollCount/%d", URL, agent.pollCount)
	_, err := client.client.R().SetHeader("Content-Type", "text/plain").Post(urlString)
	if err != nil {
		fmt.Printf("Error sending request: %s\n", err)
		return err
	}

	// reset the counter
	agent.pollCount = 0
	return err
}
