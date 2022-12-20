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
			fmt.Println("Sending data")
			m.PollCount = 0
		}
	}
}

func main() {
	newMonitor(2, 10)
}
