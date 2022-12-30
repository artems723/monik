package storage

import (
	"fmt"
)

type MemStorage struct {
	storage map[string]map[string]string
}

func NewMemStorage() *MemStorage {
	storage := make(map[string]map[string]string)
	return &MemStorage{storage: storage}
}

func (m *MemStorage) GetMetric(agentID, metricName string) (string, bool) {
	currentVal, ok := m.storage[agentID][metricName]
	return currentVal, ok
}

func (m *MemStorage) WriteMetric(agentID, metricName, metricValue string) {
	// check if agent exists in storage
	_, ok := m.storage[agentID]
	if !ok {
		// create map for agent
		m.storage[agentID] = make(map[string]string)
	}
	// add metric to storage
	m.storage[agentID][metricName] = metricValue
	fmt.Println(m.storage)
}

func (m *MemStorage) GetAllMetrics(agentID string) (map[string]string, bool) {
	allMetrics, ok := m.storage[agentID]
	return allMetrics, ok
}
