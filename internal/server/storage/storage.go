package storage

import (
	"errors"
	"github.com/artems723/monik/internal/server/domain"
)

type Repository interface {
	GetMetric(metricName string) (*domain.Metric, error)
	WriteMetric(metric *domain.Metric) error
	GetAllMetrics() (*domain.Metrics, error)
}

var ErrNotFound = errors.New("not found")
