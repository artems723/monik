package storage

import "fmt"

type Repository interface {
	Write(key, value string)
}
type MemStorage struct {
	storage map[string]string
}

func (m *MemStorage) Write(key, value string) {
	m.storage[key] = value
	fmt.Println(m.storage)
}

func NewMemStorage() *MemStorage {
	storage := make(map[string]string)
	return &MemStorage{storage: storage}
}
