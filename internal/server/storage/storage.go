package storage

import (
	"errors"
	"github.com/artems723/monik/internal/server/domain"
)

type Repository interface {
	GetMetric(metricName string) (*domain.Metrics, error)
	WriteMetric(metric *domain.Metrics) error
	GetAllMetrics() (map[string]*domain.Metrics, error)
}

var ErrNotFound = errors.New("not found")
