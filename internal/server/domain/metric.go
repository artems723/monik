package domain

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
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
	ID    string     `json:"id" db:"name"`               // имя метрики
	MType MetricType `json:"type" db:"type"`             // параметр, принимающий значение gauge или counter
	Delta *int64     `json:"delta,omitempty" db:"delta"` // значение метрики в случае передачи counter
	Value *float64   `json:"value,omitempty" db:"value"` // значение метрики в случае передачи gauge
	Hash  string     `json:"hash,omitempty" db:"-"`      // значение хеш-функции
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

func (m *Metric) String() string {
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
		*t = MetricTypeUnknown
	}
	return nil
}

func (m *Metric) Validate(key string) error {
	switch m.MType {
	case MetricTypeGauge:
		// Check that value exists
		if m.Value == nil {
			return ErrNoValue
		}
	case MetricTypeCounter:
		// Check that delta exists
		if m.Delta == nil {
			return ErrNoValue
		}
	}
	// calculate hash and compare with provided hash
	if key != "" {
		var src string
		switch m.MType {
		case MetricTypeGauge:
			src = fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value)
		case MetricTypeCounter:
			src = fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta)
		}
		if hash(src, key) != m.Hash {
			return ErrWrongKey
		}
	}
	return nil
}

func (m *Metric) AddHash(key string) {
	var src string
	switch m.MType {
	case MetricTypeGauge:
		src = fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value)
	case MetricTypeCounter:
		src = fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta)
	}
	m.Hash = hash(src, key)
}

func hash(src string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}

var ErrNoValue = errors.New("no value")
var ErrWrongKey = errors.New("wrong key for hashing")
