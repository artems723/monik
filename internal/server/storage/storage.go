package storage

import (
	"errors"
	"github.com/artems723/monik/internal/server/domain"
)

type Repository interface {
	GetMetric(metricName string) (*domain.Metric, error)
	WriteMetric(metric *domain.Metric) error
	GetAllMetrics() (map[string]*domain.Metric, error)
}

var ErrNotFound = errors.New("not found")
