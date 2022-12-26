package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

// metric types
type gauge float64
type counter int64

type monitor struct {
	Alloc         gauge
	BuckHashSys   gauge
	Frees         gauge
	GCCPUFraction gauge
	GCSys         gauge
	HeapAlloc     gauge
	HeapIdle      gauge
	HeapInuse     gauge
	HeapObjects   gauge
	HeapReleased  gauge
	HeapSys       gauge
	LastGC        gauge
	Lookups       gauge
	MCacheInuse   gauge
	MCacheSys     gauge
	MSpanInuse    gauge
	MSpanSys      gauge
	Mallocs       gauge
	NextGC        gauge
	NumForcedGC   gauge
	NumGC         gauge
	OtherSys      gauge
	PauseTotalNs  gauge
	StackInuse    gauge
	StackSys      gauge
	Sys           gauge
	TotalAlloc    gauge
	PollCount     counter
	RandomValue   gauge
}

func (m *monitor) updateMonitor() {
	var rtm runtime.MemStats

	//Read memory stats
	runtime.ReadMemStats(&rtm)

	//Update metrics
	m.Alloc = gauge(rtm.Alloc)
	m.BuckHashSys = gauge(rtm.BuckHashSys)
	m.Frees = gauge(rtm.Frees)
	m.GCCPUFraction = gauge(rtm.GCCPUFraction)
	m.GCSys = gauge(rtm.GCSys)
	m.HeapAlloc = gauge(rtm.HeapAlloc)
	m.HeapIdle = gauge(rtm.HeapIdle)
	m.HeapInuse = gauge(rtm.HeapInuse)
	m.HeapObjects = gauge(rtm.HeapObjects)
	m.HeapReleased = gauge(rtm.HeapReleased)
	m.HeapSys = gauge(rtm.HeapSys)
	m.LastGC = gauge(rtm.LastGC)
	m.Lookups = gauge(rtm.Lookups)
	m.MCacheInuse = gauge(rtm.MCacheInuse)
	m.MCacheSys = gauge(rtm.MCacheSys)
	m.MSpanInuse = gauge(rtm.MSpanInuse)
	m.MSpanSys = gauge(rtm.MSpanSys)
	m.Mallocs = gauge(rtm.Mallocs)
	m.NextGC = gauge(rtm.NextGC)
	m.NumForcedGC = gauge(rtm.NumForcedGC)
	m.NumGC = gauge(rtm.NumGC)
	m.OtherSys = gauge(rtm.OtherSys)
	m.PauseTotalNs = gauge(rtm.PauseTotalNs)
	m.StackInuse = gauge(rtm.StackInuse)
	m.StackSys = gauge(rtm.StackSys)
	m.Sys = gauge(rtm.Sys)
	m.TotalAlloc = gauge(rtm.TotalAlloc)
	m.PollCount++
	m.RandomValue = gauge(rand.Float64())
}

// send metrics data to http server
func (m *monitor) sendData() {

	endpoint := "127.0.0.1"
	port := 8080

	urlList := []string{
		fmt.Sprintf("http://%s:%d/update/gauge/Alloc/%f", endpoint, port, m.Alloc),
		fmt.Sprintf("http://%s:%d/update/gauge/BuckHashSys/%f", endpoint, port, m.BuckHashSys),
		fmt.Sprintf("http://%s:%d/update/gauge/Frees/%f", endpoint, port, m.Frees),
		fmt.Sprintf("http://%s:%d/update/gauge/GCCPUFraction/%f", endpoint, port, m.GCCPUFraction),
		fmt.Sprintf("http://%s:%d/update/gauge/GCSys/%f", endpoint, port, m.GCSys),
		fmt.Sprintf("http://%s:%d/update/gauge/HeapAlloc/%f", endpoint, port, m.HeapAlloc),
		fmt.Sprintf("http://%s:%d/update/gauge/HeapIdle/%f", endpoint, port, m.HeapIdle),
		fmt.Sprintf("http://%s:%d/update/gauge/HeapInuse/%f", endpoint, port, m.HeapInuse),
		fmt.Sprintf("http://%s:%d/update/gauge/HeapObjects/%f", endpoint, port, m.HeapObjects),
		fmt.Sprintf("http://%s:%d/update/gauge/HeapReleased/%f", endpoint, port, m.HeapReleased),
		fmt.Sprintf("http://%s:%d/update/gauge/HeapSys/%f", endpoint, port, m.HeapSys),
		fmt.Sprintf("http://%s:%d/update/gauge/LastGC/%f", endpoint, port, m.LastGC),
		fmt.Sprintf("http://%s:%d/update/gauge/Lookups/%f", endpoint, port, m.Lookups),
		fmt.Sprintf("http://%s:%d/update/gauge/MCacheInuse/%f", endpoint, port, m.MCacheInuse),
		fmt.Sprintf("http://%s:%d/update/gauge/MCacheSys/%f", endpoint, port, m.MCacheSys),
		fmt.Sprintf("http://%s:%d/update/gauge/MSpanInuse/%f", endpoint, port, m.MSpanInuse),
		fmt.Sprintf("http://%s:%d/update/gauge/MSpanSys/%f", endpoint, port, m.MSpanSys),
		fmt.Sprintf("http://%s:%d/update/gauge/Mallocs/%f", endpoint, port, m.Mallocs),
		fmt.Sprintf("http://%s:%d/update/gauge/NextGC/%f", endpoint, port, m.NextGC),
		fmt.Sprintf("http://%s:%d/update/gauge/NumForcedGC/%f", endpoint, port, m.NumForcedGC),
		fmt.Sprintf("http://%s:%d/update/gauge/NumGC/%f", endpoint, port, m.NumGC),
		fmt.Sprintf("http://%s:%d/update/gauge/OtherSys/%f", endpoint, port, m.OtherSys),
		fmt.Sprintf("http://%s:%d/update/gauge/PauseTotalNs/%f", endpoint, port, m.PauseTotalNs),
		fmt.Sprintf("http://%s:%d/update/gauge/StackInuse/%f", endpoint, port, m.StackInuse),
		fmt.Sprintf("http://%s:%d/update/gauge/StackSys/%f", endpoint, port, m.StackSys),
		fmt.Sprintf("http://%s:%d/update/gauge/Sys/%f", endpoint, port, m.Sys),
		fmt.Sprintf("http://%s:%d/update/gauge/TotalAlloc/%f", endpoint, port, m.TotalAlloc),
		fmt.Sprintf("http://%s:%d/update/counter/PollCount/%d", endpoint, port, m.PollCount),
		fmt.Sprintf("http://%s:%d/update/gauge/RandomValue/%f", endpoint, port, m.RandomValue),
	}

	for _, url := range urlList {
		fmt.Printf("Sending data to %s\n", url)
		//send metric data
		resp, err := http.Post(url, "text/plain", nil)
		if err != nil {
			fmt.Printf("Error sending request: %s\n", err)
			return
		}
		fmt.Printf("The status code we got is: %d %s\n", resp.StatusCode, http.StatusText(resp.StatusCode))
		err2 := resp.Body.Close()
		if err2 != nil {
			fmt.Printf("Error closing body: %s\n", err)
			return
		}
	}
	// reset the counter
	m.PollCount = 0
}

func newMonitor(pollInterval, reportInterval int) {
	var m monitor
	pollIntervalTicker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	reportIntervalTicker := time.NewTicker(time.Duration(reportInterval) * time.Second)

	// infinite loop for polling counters and sending it to server
	for {
		select {
		case <-pollIntervalTicker.C:
			m.updateMonitor()
			fmt.Print("Got counters: ")
			fmt.Println(m)
		case <-reportIntervalTicker.C:
			m.sendData()
		}
	}
}

func main() {
	newMonitor(2, 10)
}
