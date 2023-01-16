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

func New(s storage.Repository) *Service {
	return &Service{storage: s}
}

func (s *Service) WriteMetric(metric *domain.Metrics) error {
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
		m, err := s.storage.GetMetric(metric.ID)
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
	err := s.storage.WriteMetric(metric)
	return err
}

func (s *Service) GetMetric(metric *domain.Metrics) (*domain.Metrics, error) {
	curMetric, err := s.storage.GetMetric(metric.ID)
	if curMetric != nil && curMetric.MType != metric.MType {
		return curMetric, ErrMTypeMismatch
	}
	return curMetric, err
}

func (s *Service) GetAllMetrics() (map[string]*domain.Metrics, error) {
	return s.storage.GetAllMetrics()
}

func (s *Service) WriteMetrics(metrics []*domain.Metrics) error {
	for _, metric := range metrics {
		err := s.WriteMetric(metric)
		if err != nil {
			return err
		}
	}
	return nil
}

var ErrMTypeMismatch = errors.New("metric type mismatch")
var ErrNoValue = errors.New("no value")
