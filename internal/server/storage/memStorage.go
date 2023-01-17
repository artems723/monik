package storage

import (
	"github.com/artems723/monik/internal/server/domain"
	"log"
)

type MemStorage struct {
	s map[string]*domain.Metric
}

func NewMemStorage() *MemStorage {
	return &MemStorage{s: make(map[string]*domain.Metric)}
}

func (m *MemStorage) GetMetric(metricName string) (*domain.Metric, error) {
	currentVal, ok := m.s[metricName]
	if !ok {
		return nil, ErrNotFound
	}
	return currentVal, nil
}

func (m *MemStorage) WriteMetric(metric *domain.Metric) error {
	// add metric to storage
	m.s[metric.ID] = metric
	log.Printf("Storage was updated with metric: %v", metric)
	return nil
}

func (m *MemStorage) GetAllMetrics() (*domain.Metrics, error) {
	values := make([]*domain.Metric, 0, len(m.s))

	for _, v := range m.s {
		values = append(values, v)
	}
	return &domain.Metrics{Metrics: values}, nil
}

func (m *MemStorage) WriteAllMetrics(metrics *domain.Metrics) error {
	for _, v := range metrics.Metrics {
		m.s[v.ID] = v
	}
	return nil
}
