package storage

import (
	"github.com/artems723/monik/internal/server/domain"
	"log"
)

type MemStorage struct {
	storage map[string]map[string]*domain.Metrics
}

func NewMemStorage() *MemStorage {
	storage := make(map[string]map[string]*domain.Metrics)
	return &MemStorage{storage: storage}
}

func (m *MemStorage) GetMetric(agentID, metricName string) (*domain.Metrics, error) {
	currentVal, ok := m.storage[agentID][metricName]
	if !ok {
		return nil, ErrNotFound
	}
	return currentVal, nil
}

func (m *MemStorage) WriteMetric(agentID string, metric *domain.Metrics) error {
	// check if agent exists in storage
	_, ok := m.storage[agentID]
	if !ok {
		// create map for agent
		m.storage[agentID] = make(map[string]*domain.Metrics)
	}
	// add metric to storage
	m.storage[agentID][metric.ID] = metric
	log.Printf("Storage was updated for agent %s\n", agentID)
	return nil
}

func (m *MemStorage) GetAllMetrics(agentID string) (map[string]*domain.Metrics, error) {
	allMetrics, ok := m.storage[agentID]
	if !ok {
		return allMetrics, ErrNotFound
	}
	return allMetrics, nil
}
