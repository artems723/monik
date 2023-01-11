package storage

import (
	"errors"
	"github.com/artems723/monik/internal/server/domain"
)

type Repository interface {
	GetMetric(agentID, metricName string) (domain.Metrics, error)
	WriteMetric(agentID string, metric domain.Metrics) error
	GetAllMetrics(agentID string) (map[string]domain.Metrics, error)
}

var ErrNotFound = errors.New("not found")
