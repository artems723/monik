package service

import (
	"errors"
	"github.com/artems723/monik/internal/server"
	"github.com/artems723/monik/internal/server/domain"
	"github.com/artems723/monik/internal/server/storage"
	"log"
	"time"
)

type Service struct {
	storage  storage.Repository
	fStorage storage.Repository
	config   server.Config
}

func New(s storage.Repository, c server.Config) *Service {
	return &Service{storage: s, config: c}
}

func (s *Service) WriteMetric(metric *domain.Metric) error {
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
	// Write metric to file if storeInterval == 0s
	if s.config.StoreInterval == 0*time.Second {
		err1 := s.fStorage.WriteMetric(metric)
		if err1 != nil {
			log.Printf("Error writing metric to file: %v", err1)
		}
	}
	return err
}

func (s *Service) GetMetric(metric *domain.Metric) (*domain.Metric, error) {
	curMetric, err := s.storage.GetMetric(metric.ID)
	if curMetric != nil && curMetric.MType != metric.MType {
		return curMetric, ErrMTypeMismatch
	}
	return curMetric, err
}

func (s *Service) GetAllMetrics() (*domain.Metrics, error) {
	return s.storage.GetAllMetrics()
}

func (s *Service) WriteMetrics(metrics *domain.Metrics) error {
	for _, metric := range metrics.Metrics {
		err := s.WriteMetric(metric)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) RunFileStorage(fileStorage *storage.FileStorage) {
	// Read metrics from file to storage
	s.fStorage = fileStorage
	if s.config.Restore {
		metrics, err := s.fStorage.GetAllMetrics()
		if err != nil {
			log.Printf("error occured while reading metrics from file: %v", err)
			return
		}
		err = s.storage.WriteAllMetrics(metrics)
		if err != nil {
			log.Printf("error occured while writing metrics to storage: %v", err)
			return
		}
		log.Printf("The following metrics were loaded from file: %v", metrics)
	}

	if s.config.StoreInterval > 0*time.Second {
		// infinite loop for dumping data to file
		storeIntervalTicker := time.NewTicker(s.config.StoreInterval)
		for {
			<-storeIntervalTicker.C
			err := s.WriteAllToFile()
			if err != nil {
				log.Printf("error occured while dumping data to file: %v", err)
				return
			}
			log.Printf("Stored to file")
		}
	}
}

func (s *Service) WriteAllToFile() error {
	// Read all metrics from storage
	metrics, err := s.storage.GetAllMetrics()
	if err != nil {
		log.Printf("error occured while reading all metrics from storage: %v", err)
		return err
	}
	// Write all metrics to file
	err = s.fStorage.WriteAllMetrics(metrics)
	if err != nil {
		log.Printf("error occured while dumping data to file: %v", err)
		return err
	}
	return nil
}

func (s *Service) Shutdown() error {
	err := s.WriteAllToFile()
	if err != nil {
		log.Printf("error occured while dumping data to file: %v", err)
		return err
	}
	log.Printf("Stored to file before shutdown")
	return nil
}

var ErrMTypeMismatch = errors.New("metric type mismatch")
var ErrNoValue = errors.New("no value")
