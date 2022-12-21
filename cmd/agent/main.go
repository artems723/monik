package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

type gauge float64
type counter int64

type monitor struct {
	rtm           runtime.MemStats
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
	//Read memory stats
	runtime.ReadMemStats(&m.rtm)

	//Update metrics
	m.Alloc = gauge(m.rtm.Alloc)
	m.BuckHashSys = gauge(m.rtm.BuckHashSys)
	m.Frees = gauge(m.rtm.Frees)
	m.GCCPUFraction = gauge(m.rtm.GCCPUFraction)
	m.GCSys = gauge(m.rtm.GCSys)
	m.HeapAlloc = gauge(m.rtm.HeapAlloc)
	m.HeapIdle = gauge(m.rtm.HeapIdle)
	m.HeapInuse = gauge(m.rtm.HeapInuse)
	m.HeapObjects = gauge(m.rtm.HeapObjects)
	m.HeapReleased = gauge(m.rtm.HeapReleased)
	m.HeapSys = gauge(m.rtm.HeapSys)
	m.LastGC = gauge(m.rtm.LastGC)
	m.Lookups = gauge(m.rtm.Lookups)
	m.MCacheInuse = gauge(m.rtm.MCacheInuse)
	m.MCacheSys = gauge(m.rtm.MCacheSys)
	m.MSpanInuse = gauge(m.rtm.MSpanInuse)
	m.MSpanSys = gauge(m.rtm.MSpanSys)
	m.Mallocs = gauge(m.rtm.Mallocs)
	m.NextGC = gauge(m.rtm.NextGC)
	m.NumForcedGC = gauge(m.rtm.NumForcedGC)
	m.NumGC = gauge(m.rtm.NumGC)
	m.OtherSys = gauge(m.rtm.OtherSys)
	m.PauseTotalNs = gauge(m.rtm.PauseTotalNs)
	m.StackInuse = gauge(m.rtm.StackInuse)
	m.StackSys = gauge(m.rtm.StackSys)
	m.Sys = gauge(m.rtm.Sys)
	m.TotalAlloc = gauge(m.rtm.TotalAlloc)
	m.PollCount++
	m.RandomValue = gauge(rand.Float64())
}

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
		//_, err := http.Post(url, "text/plain", nil)
		//if err != nil {
		//	// handle error
		//}

		fmt.Printf("Sending data to %s\n", url)
	}

	m.PollCount = 0
}

func newMonitor(pollInterval, reportInterval int) {
	var m monitor

	pollIntervalTicker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	reportIntervalTicker := time.NewTicker(time.Duration(reportInterval) * time.Second)
	for {
		select {
		case <-pollIntervalTicker.C:

			m.updateMonitor()
			fmt.Println(m)

		case <-reportIntervalTicker.C:
			m.sendData()
		}
	}
}

func main() {
	newMonitor(2, 10)
}
