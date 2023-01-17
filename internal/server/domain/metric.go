package domain

import (
	"errors"
	"fmt"
	"strings"
)

type MetricType string

const (
	MetricTypeGauge   MetricType = "gauge"
	MetricTypeCounter MetricType = "counter"
)

type Metric struct {
	ID    string     `json:"id"`              // имя метрики
	MType MetricType `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64     `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64   `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type Metrics struct {
	Metrics []*Metric `json:"metrics"`
}

func NewGaugeMetric(id string, value float64) *Metric {
	return &Metric{ID: id, MType: MetricTypeGauge, Value: &value}
}

func NewCounterMetric(id string, delta int64) *Metric {
	return &Metric{ID: id, MType: MetricTypeCounter, Delta: &delta}
}

func NewMetric(id, mType string) *Metric {
	return &Metric{ID: id, MType: MetricType(mType)}
}

func (m Metric) String() string {
	// check metric type
	switch m.MType {
	case MetricTypeGauge:
		if m.Value != nil {
			return fmt.Sprintf("Name: %s, Type: %s, Value: %f", m.ID, m.MType, *m.Value)
		}
	case MetricTypeCounter:
		if m.Delta != nil {
			return fmt.Sprintf("Name: %s, Type: %s, Delta: %d", m.ID, m.MType, *m.Delta)
		}
	}
	return fmt.Sprintf("Name: %s, Type: %s", m.ID, m.MType)
}

func (t *MetricType) UnmarshalJSON(data []byte) error {
	m := MetricType(strings.Trim(string(data), "\""))
	switch m {
	case MetricTypeGauge, MetricTypeCounter:
		*t = m
	default:
		return ErrUnknownMetricType
	}
	return nil
}

var ErrUnknownMetricType = errors.New("unknown metric type")
