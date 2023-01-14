package client

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
)

// metric types
type (
	metricTypeGauge   float64
	metricTypeCounter int64

	Agent struct {
		id           string
		gaugeMetrics map[string]metricTypeGauge
		pollCount    metricTypeCounter
	}
)

func NewAgent() Agent {
	metricsMap := make(map[string]metricTypeGauge)
	return Agent{gaugeMetrics: metricsMap, pollCount: 0}
}

func (agent *Agent) UpdateMetrics() {
	var rtm runtime.MemStats

	//Read memory stats
	runtime.ReadMemStats(&rtm)

	//Update metrics
	agent.gaugeMetrics["Alloc"] = metricTypeGauge(rtm.Alloc)
	agent.gaugeMetrics["BuckHashSys"] = metricTypeGauge(rtm.BuckHashSys)
	agent.gaugeMetrics["Frees"] = metricTypeGauge(rtm.Frees)
	agent.gaugeMetrics["GCCPUFraction"] = metricTypeGauge(rtm.GCCPUFraction)
	agent.gaugeMetrics["GCSys"] = metricTypeGauge(rtm.GCSys)
	agent.gaugeMetrics["HeapAlloc"] = metricTypeGauge(rtm.HeapAlloc)
	agent.gaugeMetrics["HeapIdle"] = metricTypeGauge(rtm.HeapIdle)
	agent.gaugeMetrics["HeapInuse"] = metricTypeGauge(rtm.HeapInuse)
	agent.gaugeMetrics["HeapObjects"] = metricTypeGauge(rtm.HeapObjects)
	agent.gaugeMetrics["HeapReleased"] = metricTypeGauge(rtm.HeapReleased)
	agent.gaugeMetrics["HeapSys"] = metricTypeGauge(rtm.HeapSys)
	agent.gaugeMetrics["LastGC"] = metricTypeGauge(rtm.LastGC)
	agent.gaugeMetrics["Lookups"] = metricTypeGauge(rtm.Lookups)
	agent.gaugeMetrics["MCacheInuse"] = metricTypeGauge(rtm.MCacheInuse)
	agent.gaugeMetrics["MCacheSys"] = metricTypeGauge(rtm.MCacheSys)
	agent.gaugeMetrics["MSpanInuse"] = metricTypeGauge(rtm.MSpanInuse)
	agent.gaugeMetrics["MSpanSys"] = metricTypeGauge(rtm.MSpanSys)
	agent.gaugeMetrics["Mallocs"] = metricTypeGauge(rtm.Mallocs)
	agent.gaugeMetrics["NextGC"] = metricTypeGauge(rtm.NextGC)
	agent.gaugeMetrics["NumForcedGC"] = metricTypeGauge(rtm.NumForcedGC)
	agent.gaugeMetrics["NumGC"] = metricTypeGauge(rtm.NumGC)
	agent.gaugeMetrics["OtherSys"] = metricTypeGauge(rtm.OtherSys)
	agent.gaugeMetrics["PauseTotalNs"] = metricTypeGauge(rtm.PauseTotalNs)
	agent.gaugeMetrics["StackInuse"] = metricTypeGauge(rtm.StackInuse)
	agent.gaugeMetrics["StackSys"] = metricTypeGauge(rtm.StackSys)
	agent.gaugeMetrics["Sys"] = metricTypeGauge(rtm.Sys)
	agent.gaugeMetrics["TotalAlloc"] = metricTypeGauge(rtm.TotalAlloc)
	agent.pollCount++
	agent.gaugeMetrics["RandomValue"] = metricTypeGauge(rand.Float64())
}

// send metrics data to http server
func (agent *Agent) SendData(URL string, client HTTPClient) {
	log.Printf("Sending data")
	// send gauges
	for key, val := range agent.gaugeMetrics {
		urlString := fmt.Sprintf("%s/update/gauge/%s/%f", URL, key, val)

		//log.Printf("Sending data to %s\n", urlString)
		_, err := client.client.R().SetHeader("Content-Type", "text/plain").Post(urlString)
		if err != nil {
			log.Printf("Error sending request: %s\n", err)
			return
		}
	}

	// send counter
	urlString := fmt.Sprintf("%s/update/counter/PollCount/%d", URL, agent.pollCount)
	_, err := client.client.R().SetHeader("Content-Type", "text/plain").Post(urlString)
	if err != nil {
		log.Printf("Error sending request: %s\n", err)
		return
	}

	// reset the counter
	agent.pollCount = 0
}
