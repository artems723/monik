package agent

import (
	"fmt"
	"github.com/artems723/monik/internal/client"
	"github.com/artems723/monik/internal/client/httpClient"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"log"
	"math/rand"
	"runtime"
	"sync"
)

type Agent struct {
	storage map[string]*client.Metric
	client  httpClient.Client
	key     string
	mu      *sync.RWMutex
}

func NewAgent(key string, cl httpClient.Client) Agent {
	return Agent{
		storage: make(map[string]*client.Metric),
		client:  cl,
		key:     key,
		mu:      &sync.RWMutex{},
	}
}

func (agent *Agent) UpdateMetrics() {
	var rtm runtime.MemStats
	// Read memory stats
	runtime.ReadMemStats(&rtm)
	// Update metrics
	agent.mu.Lock()
	defer agent.mu.Unlock()
	agent.storage["Alloc"] = client.NewGaugeMetric("Alloc", float64(rtm.Alloc), agent.key)
	agent.storage["BuckHashSys"] = client.NewGaugeMetric("BuckHashSys", float64(rtm.BuckHashSys), agent.key)
	agent.storage["Frees"] = client.NewGaugeMetric("Frees", float64(rtm.Frees), agent.key)
	agent.storage["GCCPUFraction"] = client.NewGaugeMetric("GCCPUFraction", rtm.GCCPUFraction, agent.key)
	agent.storage["GCSys"] = client.NewGaugeMetric("GCSys", float64(rtm.GCSys), agent.key)
	agent.storage["HeapAlloc"] = client.NewGaugeMetric("HeapAlloc", float64(rtm.HeapAlloc), agent.key)
	agent.storage["HeapIdle"] = client.NewGaugeMetric("HeapIdle", float64(rtm.HeapIdle), agent.key)
	agent.storage["HeapInuse"] = client.NewGaugeMetric("HeapInuse", float64(rtm.HeapInuse), agent.key)
	agent.storage["HeapObjects"] = client.NewGaugeMetric("HeapObjects", float64(rtm.HeapObjects), agent.key)
	agent.storage["HeapReleased"] = client.NewGaugeMetric("HeapReleased", float64(rtm.HeapReleased), agent.key)
	agent.storage["HeapSys"] = client.NewGaugeMetric("HeapSys", float64(rtm.HeapSys), agent.key)
	agent.storage["LastGC"] = client.NewGaugeMetric("LastGC", float64(rtm.LastGC), agent.key)
	agent.storage["Lookups"] = client.NewGaugeMetric("Lookups", float64(rtm.Lookups), agent.key)
	agent.storage["MCacheInuse"] = client.NewGaugeMetric("MCacheInuse", float64(rtm.MCacheInuse), agent.key)
	agent.storage["MCacheSys"] = client.NewGaugeMetric("MCacheSys", float64(rtm.MCacheSys), agent.key)
	agent.storage["MSpanInuse"] = client.NewGaugeMetric("MSpanInuse", float64(rtm.MSpanInuse), agent.key)
	agent.storage["MSpanSys"] = client.NewGaugeMetric("MSpanSys", float64(rtm.MSpanSys), agent.key)
	agent.storage["Mallocs"] = client.NewGaugeMetric("Mallocs", float64(rtm.Mallocs), agent.key)
	agent.storage["NextGC"] = client.NewGaugeMetric("NextGC", float64(rtm.NextGC), agent.key)
	agent.storage["NumForcedGC"] = client.NewGaugeMetric("NumForcedGC", float64(rtm.NumForcedGC), agent.key)
	agent.storage["NumGC"] = client.NewGaugeMetric("NumGC", float64(rtm.NumGC), agent.key)
	agent.storage["OtherSys"] = client.NewGaugeMetric("OtherSys", float64(rtm.OtherSys), agent.key)
	agent.storage["PauseTotalNs"] = client.NewGaugeMetric("PauseTotalNs", float64(rtm.PauseTotalNs), agent.key)
	agent.storage["StackInuse"] = client.NewGaugeMetric("StackInuse", float64(rtm.StackInuse), agent.key)
	agent.storage["StackSys"] = client.NewGaugeMetric("StackSys", float64(rtm.StackSys), agent.key)
	agent.storage["Sys"] = client.NewGaugeMetric("Sys", float64(rtm.Sys), agent.key)
	agent.storage["TotalAlloc"] = client.NewGaugeMetric("TotalAlloc", float64(rtm.TotalAlloc), agent.key)
	agent.storage["RandomValue"] = client.NewGaugeMetric("RandomValue", rand.Float64(), agent.key)
	// Update counter
	var delta int64 = 1
	m, ok := agent.storage["PollCount"]
	if ok {
		delta += *m.Delta
	}
	agent.storage["PollCount"] = client.NewCounterMetric("PollCount", delta, agent.key)
	log.Printf("Got counters. PollCount=%d", *agent.storage["PollCount"].Delta)
}

func (agent *Agent) getValues() []*client.Metric {
	agent.mu.RLock()
	defer agent.mu.RUnlock()
	values := make([]*client.Metric, 0, len(agent.storage))
	for _, v := range agent.storage {
		values = append(values, v)
	}
	return values
}

// Send metrics to http server
func (agent *Agent) SendData(serverAddr string) {
	URL := fmt.Sprintf("%s/updates/", serverAddr)
	metrics := agent.getValues()
	m, err := agent.client.SendData(metrics, URL)
	if err != nil {
		return
	}
	log.Printf("Got response from server: %v", m)
	agent.resetCounter()
	log.Printf("Metrics were succesfully sent")
}

func (agent *Agent) resetCounter() {
	agent.mu.Lock()
	defer agent.mu.Unlock()
	if _, ok := agent.storage["PollCount"]; ok {
		*agent.storage["PollCount"].Delta = 0
	}
}

func (agent *Agent) UpdateAdditionalMetrics() {
	v, _ := mem.VirtualMemory()
	c, _ := cpu.Percent(0, false)
	// Update metrics
	agent.mu.Lock()
	defer agent.mu.Unlock()
	agent.storage["TotalMemory"] = client.NewGaugeMetric("TotalMemory", float64(v.Total), agent.key)
	agent.storage["FreeMemory"] = client.NewGaugeMetric("FreeMemory", float64(v.Free), agent.key)
	agent.storage["CPUutilization1"] = client.NewGaugeMetric("CPUutilization1", c[0], agent.key)
	log.Printf("Got additional counters")
}
