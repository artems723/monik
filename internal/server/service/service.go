package service

import (
	"errors"
	"github.com/artems723/monik/internal/server/domain"
	"github.com/artems723/monik/internal/server/storage"
	"log"
)

type Service struct {
	storage storage.Repository
}

func New(s storage.Repository) Service {
	return Service{storage: s}
}

func (s Service) WriteMetric(agentID string, metric domain.Metrics) error {
	// Check metric type
	switch metric.MType {
	case domain.MetricTypeGauge:
		// Check that value exists
		if metric.Value == nil {
			return ErrNoValue
		}
	case domain.MetricTypeCounter:
		// Check that delta exists
		if metric.Delta == nil {
			return ErrNoValue
		}
		// Get current metric from storage to sum deltas
		m, err := s.storage.GetMetric(agentID, metric.ID)
		// Check for errors
		if err != nil && !errors.Is(err, storage.ErrNotFound) {
			log.Printf("storage.GetMetric: %v", err)
			return err
		}
		if errors.Is(err, storage.ErrNotFound) {
			break
		}
		// Check if current metric in not 'counter' type
		if m.MType != domain.MetricTypeCounter {
			return ErrMTypeMismatch
		}
		// Add delta to current value
		*metric.Delta += *m.Delta
	default:
		return domain.ErrUnknownMetricType
	}
	// Write metric to storage
	err := s.storage.WriteMetric(agentID, metric)
	return err
}

func (s Service) GetMetric(agentID, metricName string) (domain.Metrics, error) {
	return s.storage.GetMetric(agentID, metricName)
}

func (s Service) GetAllMetrics(agentID string) (map[string]domain.Metrics, error) {
	return s.storage.GetAllMetrics(agentID)
}

var ErrMTypeMismatch = errors.New("metric type mismatch")
var ErrNoValue = errors.New("no value")
