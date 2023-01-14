package client

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"runtime"
)

type Agent struct {
	storage map[string]Metrics
}

func NewAgent() Agent {
	m := make(map[string]Metrics)
	m["pollCount"] = NewCounterMetric("pollCount", 0)
	return Agent{storage: m}
}

func (agent *Agent) UpdateMetrics() {
	var rtm runtime.MemStats

	//Read memory stats
	runtime.ReadMemStats(&rtm)

	//Update metrics
	agent.storage["Alloc"] = NewGaugeMetric("Alloc", float64(rtm.Alloc))
	agent.storage["BuckHashSys"] = NewGaugeMetric("BuckHashSys", float64(rtm.BuckHashSys))
	agent.storage["Frees"] = NewGaugeMetric("Frees", float64(rtm.Frees))
	agent.storage["GCCPUFraction"] = NewGaugeMetric("GCCPUFraction", rtm.GCCPUFraction)
	agent.storage["GCSys"] = NewGaugeMetric("GCSys", float64(rtm.GCSys))
	agent.storage["HeapAlloc"] = NewGaugeMetric("HeapAlloc", float64(rtm.HeapAlloc))
	agent.storage["HeapIdle"] = NewGaugeMetric("HeapIdle", float64(rtm.HeapIdle))
	agent.storage["HeapInuse"] = NewGaugeMetric("HeapInuse", float64(rtm.HeapInuse))
	agent.storage["HeapObjects"] = NewGaugeMetric("HeapObjects", float64(rtm.HeapObjects))
	agent.storage["HeapReleased"] = NewGaugeMetric("HeapReleased", float64(rtm.HeapReleased))
	agent.storage["HeapSys"] = NewGaugeMetric("HeapSys", float64(rtm.HeapSys))
	agent.storage["LastGC"] = NewGaugeMetric("LastGC", float64(rtm.LastGC))
	agent.storage["Lookups"] = NewGaugeMetric("Lookups", float64(rtm.Lookups))
	agent.storage["MCacheInuse"] = NewGaugeMetric("MCacheInuse", float64(rtm.MCacheInuse))
	agent.storage["MCacheSys"] = NewGaugeMetric("MCacheSys", float64(rtm.MCacheSys))
	agent.storage["MSpanInuse"] = NewGaugeMetric("MSpanInuse", float64(rtm.MSpanInuse))
	agent.storage["MSpanSys"] = NewGaugeMetric("MSpanSys", float64(rtm.MSpanSys))
	agent.storage["Mallocs"] = NewGaugeMetric("Mallocs", float64(rtm.Mallocs))
	agent.storage["NextGC"] = NewGaugeMetric("NextGC", float64(rtm.NextGC))
	agent.storage["NumForcedGC"] = NewGaugeMetric("NumForcedGC", float64(rtm.NumForcedGC))
	agent.storage["NumGC"] = NewGaugeMetric("NumGC", float64(rtm.NumGC))
	agent.storage["OtherSys"] = NewGaugeMetric("OtherSys", float64(rtm.OtherSys))
	agent.storage["PauseTotalNs"] = NewGaugeMetric("PauseTotalNs", float64(rtm.PauseTotalNs))
	agent.storage["StackInuse"] = NewGaugeMetric("StackInuse", float64(rtm.StackInuse))
	agent.storage["StackSys"] = NewGaugeMetric("StackSys", float64(rtm.StackSys))
	agent.storage["Sys"] = NewGaugeMetric("Sys", float64(rtm.Sys))
	agent.storage["TotalAlloc"] = NewGaugeMetric("TotalAlloc", float64(rtm.TotalAlloc))
	*agent.storage["pollCount"].Delta++
	agent.storage["RandomValue"] = NewGaugeMetric("RandomValue", rand.Float64())
	log.Printf("Got counters")
}

// send metrics data to http server
func (agent *Agent) SendData(URL string, client HTTPClient) {
	// send metrics
	for _, metric := range agent.storage {
		urlString := fmt.Sprintf("%s/update/", URL)

		m, err := json.Marshal(metric)
		if err != nil {
			log.Printf("agent.SendData: unable to marshal. Error: %v. Metric: %v", err, metric)
			return
		}
		_, err = client.client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(m).
			Post(urlString)
		if err != nil {
			log.Printf("Error sending request: %s", err)
			return
		}
	}
	// reset the counter
	*agent.storage["pollCount"].Delta = 0
	log.Printf("Metrics were succesfully sent")
}
