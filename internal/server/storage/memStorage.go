package storage

import (
	"github.com/artems723/monik/internal/server/domain"
	"log"
)

type MemStorage struct {
	Storage map[string]*domain.Metric
}

func NewMemStorage() MemStorage {
	storage := make(map[string]*domain.Metric)
	return MemStorage{Storage: storage}
}

func (m *MemStorage) GetMetric(metricName string) (*domain.Metric, error) {
	currentVal, ok := m.Storage[metricName]
	if !ok {
		return nil, ErrNotFound
	}
	return currentVal, nil
}

func (m *MemStorage) WriteMetric(metric *domain.Metric) error {
	// add metric to storage
	m.Storage[metric.ID] = metric
	log.Printf("Storage was updated with metric: %v", metric)
	return nil
}

func (m *MemStorage) GetAllMetrics() (*domain.Metrics, error) {
	values := make([]*domain.Metric, 0, len(m.Storage))

	for _, v := range m.Storage {
		values = append(values, v)
	}
	return &domain.Metrics{Metrics: values}, nil
}
