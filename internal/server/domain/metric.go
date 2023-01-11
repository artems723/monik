package domain

import "fmt"

type MetricType string

const (
	MetricTypeGauge   MetricType = "gauge"
	MetricTypeCounter MetricType = "counter"
)

type Metrics struct {
	ID    string     `json:"id"`              // имя метрики
	MType MetricType `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64     `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64   `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func NewMetric(id string, mtype MetricType, delta int64, value float64) Metrics {
	return Metrics{id, mtype, &delta, &value}
}

func NewGaugeMetric(id string, value float64) Metrics {
	return Metrics{ID: id, MType: MetricTypeGauge, Value: &value}
}

func NewCounterMetric(id string, delta int64) Metrics {
	return Metrics{ID: id, MType: MetricTypeCounter, Delta: &delta}
}

func (m Metrics) String() string {
	var s string
	switch m.MType {
	case MetricTypeGauge:
		s = fmt.Sprintf("ID: %s, Mtype: %s, Value: %f", m.ID, m.MType, *m.Value)
	case MetricTypeCounter:
		s = fmt.Sprintf("ID: %s, Mtype: %s, Delta: %d", m.ID, m.MType, *m.Delta)
	}
	return s
}
