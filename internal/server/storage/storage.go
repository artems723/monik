package storage

import (
	"fmt"
	"math/big"
)

type Repository interface {
	Write(metricType, key, value string)
}
type MemStorage struct {
	storage map[string]string
}

func (m *MemStorage) Write(metricType, key, value string) {
	switch metricType {
	case "counter":
		var big1, big2 *big.Int
		currentVal, ok := m.storage[key]
		if !ok {
			big1 = big.NewInt(0)
		} else {
			big1, _ = new(big.Int).SetString(currentVal, 0)
		}
		big2, ok = new(big.Int).SetString(value, 0)
		if !ok {
			fmt.Printf("%s is not an integer.", value)
		} else {
			big1.Add(big1, big2)
			str := fmt.Sprintf("%v", big1)
			m.storage[key] = str
		}
	default:
		m.storage[key] = value
	}

	fmt.Println(m.storage)
}

func NewMemStorage() *MemStorage {
	storage := make(map[string]string)
	return &MemStorage{storage: storage}
}
