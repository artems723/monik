package client

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type MetricType string

const (
	MetricTypeGauge   MetricType = "gauge"
	MetricTypeCounter MetricType = "counter"
	MetricTypeUnknown MetricType = "unknown"
)

type Metric struct {
	ID    string     `json:"id"`              // имя метрики
	MType MetricType `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64     `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64   `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string     `json:"hash,omitempty"`  // значение хеш-функции
}

func NewGaugeMetric(id string, value float64, key string) *Metric {
	// calculate hash if key exists
	if key != "" {
		return &Metric{
			ID:    id,
			MType: MetricTypeGauge,
			Value: &value,
			Hash:  hash(fmt.Sprintf("%s:gauge:%f", id, value), key),
		}
	}
	return &Metric{
		ID:    id,
		MType: MetricTypeGauge,
		Value: &value,
	}
}

func NewCounterMetric(id string, delta int64, key string) *Metric {
	// calculate hash if key exists
	if key != "" {
		return &Metric{
			ID:    id,
			MType: MetricTypeCounter,
			Delta: &delta,
			Hash:  hash(fmt.Sprintf("%s:counter:%d", id, delta), key),
		}
	}
	return &Metric{
		ID:    id,
		MType: MetricTypeCounter,
		Delta: &delta,
	}
}

func hash(src string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}

func (m *Metric) String() string {
	// check metric type
	switch m.MType {
	case MetricTypeGauge:
		if m.Value != nil {
			return fmt.Sprintf("ID: %s, Mtype: %s, Value: %f", m.ID, m.MType, *m.Value)
		}
	case MetricTypeCounter:
		if m.Delta != nil {
			return fmt.Sprintf("ID: %s, Mtype: %s, Delta: %d", m.ID, m.MType, *m.Delta)
		}
	}
	return fmt.Sprintf("ID: %s, Mtype: %s", m.ID, m.MType)
}

func (t *MetricType) UnmarshalJSON(data []byte) error {
	m := MetricType(strings.Trim(string(data), "\""))
	switch m {
	case MetricTypeGauge, MetricTypeCounter:
		*t = m
	default:
		*t = MetricTypeUnknown
	}
	return nil
}
