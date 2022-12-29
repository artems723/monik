package storage

import (
	"fmt"
)

type MemStorage struct {
	storage map[string]string
}

func NewMemStorage() *MemStorage {
	storage := make(map[string]string)
	return &MemStorage{storage: storage}
}

func (m *MemStorage) Get(metricName string) (string, bool) {
	currentVal, ok := m.storage[metricName]
	return currentVal, ok
}

func (m *MemStorage) Write(metricName, metricValue string) {
	m.storage[metricName] = metricValue
	fmt.Println(m.storage)
}
