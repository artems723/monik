package client

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"runtime"
)

type Agent struct {
	storage map[string]*Metrics
}

func NewAgent() Agent {
	metricsMap := make(map[string]*Metrics)

	metricsMap["Alloc"] = NewMetric("Alloc", MetricTypeGauge)
	metricsMap["BuckHashSys"] = NewMetric("BuckHashSys", MetricTypeGauge)
	metricsMap["Frees"] = NewMetric("Frees", MetricTypeGauge)
	metricsMap["GCCPUFraction"] = NewMetric("GCCPUFraction", MetricTypeGauge)
	metricsMap["GCSys"] = NewMetric("GCSys", MetricTypeGauge)
	metricsMap["HeapAlloc"] = NewMetric("HeapAlloc", MetricTypeGauge)
	metricsMap["HeapIdle"] = NewMetric("HeapIdle", MetricTypeGauge)
	metricsMap["HeapInuse"] = NewMetric("HeapInuse", MetricTypeGauge)
	metricsMap["HeapObjects"] = NewMetric("HeapObjects", MetricTypeGauge)
	metricsMap["HeapReleased"] = NewMetric("HeapReleased", MetricTypeGauge)
	metricsMap["HeapSys"] = NewMetric("HeapSys", MetricTypeGauge)
	metricsMap["LastGC"] = NewMetric("LastGC", MetricTypeGauge)
	metricsMap["Lookups"] = NewMetric("Lookups", MetricTypeGauge)
	metricsMap["MCacheInuse"] = NewMetric("MCacheInuse", MetricTypeGauge)
	metricsMap["MCacheSys"] = NewMetric("MCacheSys", MetricTypeGauge)
	metricsMap["MSpanInuse"] = NewMetric("MSpanInuse", MetricTypeGauge)
	metricsMap["MSpanSys"] = NewMetric("MSpanSys", MetricTypeGauge)
	metricsMap["Mallocs"] = NewMetric("Mallocs", MetricTypeGauge)
	metricsMap["NextGC"] = NewMetric("NextGC", MetricTypeGauge)
	metricsMap["NumForcedGC"] = NewMetric("NumForcedGC", MetricTypeGauge)
	metricsMap["NumGC"] = NewMetric("NumGC", MetricTypeGauge)
	metricsMap["OtherSys"] = NewMetric("OtherSys", MetricTypeGauge)
	metricsMap["PauseTotalNs"] = NewMetric("PauseTotalNs", MetricTypeGauge)
	metricsMap["StackInuse"] = NewMetric("StackInuse", MetricTypeGauge)
	metricsMap["StackSys"] = NewMetric("StackSys", MetricTypeGauge)
	metricsMap["Sys"] = NewMetric("Sys", MetricTypeGauge)
	metricsMap["TotalAlloc"] = NewMetric("TotalAlloc", MetricTypeGauge)
	metricsMap["pollCount"] = NewCounterMetric("pollCount", 0)
	metricsMap["RandomValue"] = NewMetric("RandomValue", MetricTypeGauge)

	return Agent{storage: metricsMap}
}

func (agent *Agent) UpdateMetrics() {
	var rtm runtime.MemStats

	//Read memory stats
	runtime.ReadMemStats(&rtm)

	//Update metrics
	agent.storage["Alloc"].Value = createFloat64(float64(rtm.Alloc))
	agent.storage["BuckHashSys"].Value = createFloat64(float64(rtm.BuckHashSys))
	agent.storage["Frees"].Value = createFloat64(float64(rtm.Frees))
	agent.storage["GCCPUFraction"].Value = createFloat64(rtm.GCCPUFraction)
	agent.storage["GCSys"].Value = createFloat64(float64(rtm.GCSys))
	agent.storage["HeapAlloc"].Value = createFloat64(float64(rtm.HeapAlloc))
	agent.storage["HeapIdle"].Value = createFloat64(float64(rtm.HeapIdle))
	agent.storage["HeapInuse"].Value = createFloat64(float64(rtm.HeapInuse))
	agent.storage["HeapObjects"].Value = createFloat64(float64(rtm.HeapObjects))
	agent.storage["HeapReleased"].Value = createFloat64(float64(rtm.HeapReleased))
	agent.storage["HeapSys"].Value = createFloat64(float64(rtm.HeapSys))
	agent.storage["LastGC"].Value = createFloat64(float64(rtm.LastGC))
	agent.storage["Lookups"].Value = createFloat64(float64(rtm.Lookups))
	agent.storage["MCacheInuse"].Value = createFloat64(float64(rtm.MCacheInuse))
	agent.storage["MCacheSys"].Value = createFloat64(float64(rtm.MCacheSys))
	agent.storage["MSpanInuse"].Value = createFloat64(float64(rtm.MSpanInuse))
	agent.storage["MSpanSys"].Value = createFloat64(float64(rtm.MSpanSys))
	agent.storage["Mallocs"].Value = createFloat64(float64(rtm.Mallocs))
	agent.storage["NextGC"].Value = createFloat64(float64(rtm.NextGC))
	agent.storage["NumForcedGC"].Value = createFloat64(float64(rtm.NumForcedGC))
	agent.storage["NumGC"].Value = createFloat64(float64(rtm.NumGC))
	agent.storage["OtherSys"].Value = createFloat64(float64(rtm.OtherSys))
	agent.storage["PauseTotalNs"].Value = createFloat64(float64(rtm.PauseTotalNs))
	agent.storage["StackInuse"].Value = createFloat64(float64(rtm.StackInuse))
	agent.storage["StackSys"].Value = createFloat64(float64(rtm.StackSys))
	agent.storage["Sys"].Value = createFloat64(float64(rtm.Sys))
	agent.storage["TotalAlloc"].Value = createFloat64(float64(rtm.TotalAlloc))
	*agent.storage["pollCount"].Delta++
	agent.storage["RandomValue"].Value = createFloat64(rand.Float64())
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

func createFloat64(x float64) *float64 {
	return &x
}
