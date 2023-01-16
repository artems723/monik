package storage

import (
	"github.com/artems723/monik/internal/server/domain"
	"log"
)

type MemStorage struct {
	storage map[string]*domain.Metrics
}

func NewMemStorage() *MemStorage {
	storage := make(map[string]*domain.Metrics)
	return &MemStorage{storage: storage}
}

func (m *MemStorage) GetMetric(metricName string) (*domain.Metrics, error) {
	currentVal, ok := m.storage[metricName]
	if !ok {
		return nil, ErrNotFound
	}
	return currentVal, nil
}

func (m *MemStorage) WriteMetric(metric *domain.Metrics) error {
	// add metric to storage
	m.storage[metric.ID] = metric
	log.Printf("Storage was updated with metric: %v", metric)
	return nil
}

func (m *MemStorage) GetAllMetrics() (map[string]*domain.Metrics, error) {
	return m.storage, nil
}
