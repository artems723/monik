package storage

import (
	"errors"
	"github.com/artems723/monik/internal/server/domain"
	"golang.org/x/net/context"
)

type Repository interface {
	GetMetric(ctx context.Context, metricName string) (*domain.Metric, error)
	WriteMetric(ctx context.Context, metric *domain.Metric) error
	GetAllMetrics(ctx context.Context) (*domain.Metrics, error)
	WriteAllMetrics(ctx context.Context, metrics *domain.Metrics) error
	PingRepo() error
}

var ErrNotFound = errors.New("not found")
