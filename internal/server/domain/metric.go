package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

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

func NewGaugeMetric(id string, value float64) Metrics {
	return Metrics{ID: id, MType: MetricTypeGauge, Value: &value}
}

func NewCounterMetric(id string, delta int64) Metrics {
	return Metrics{ID: id, MType: MetricTypeCounter, Delta: &delta}
}

func (m Metrics) String() string {
	var s string
	// check metric type
	switch m.MType {
	case MetricTypeGauge:
		s = fmt.Sprintf("ID: %s, Mtype: %s, Value: %f", m.ID, m.MType, *m.Value)
	case MetricTypeCounter:
		s = fmt.Sprintf("ID: %s, Mtype: %s, Delta: %d", m.ID, m.MType, *m.Delta)
	}
	return s
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

func (t *Metrics) UnmarshalJSON(data []byte) error {
	type MetricsAlias Metrics
	var aliasValue MetricsAlias
	err := json.Unmarshal(data, &aliasValue)
	if err != nil {
		return err
	}
	switch aliasValue.MType {
	case MetricTypeGauge:
		if aliasValue.Value == nil {
			return ErrNoValue
		}
	case MetricTypeCounter:
		if aliasValue.Delta == nil {
			return ErrNoValue
		}
	default:
		return ErrUnknownMetricType
	}
	*t = Metrics(aliasValue)
	return nil
}

var ErrUnknownMetricType = errors.New("unknown metric type")
var ErrNoValue = errors.New("no value")
