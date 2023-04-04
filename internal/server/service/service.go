// Package service contains all business logic for server
package service

import (
	"context"
	"log"
	"time"

	"github.com/artems723/monik/internal/server/config"
	"github.com/artems723/monik/internal/server/domain"
	"github.com/artems723/monik/internal/server/storage"
	"github.com/pkg/errors"
)

type Service struct {
	storage  Repository
	fStorage *storage.FileStorage
	config   config.Config
}

func New(s Repository, c config.Config) *Service {
	return &Service{storage: s, config: c}
}

func (s *Service) WriteMetric(ctx context.Context, metric *domain.Metric) (*domain.Metric, error) {
	// Increment delta value of counter metric
	if metric.MType == domain.MetricTypeCounter {
		// Get current metric from storage to sum deltas
		m, err := s.storage.GetMetric(ctx, metric.ID)
		// Check for errors
		if err != nil && !errors.Is(err, storage.ErrNotFound) {
			return nil, errors.New("storage.GetMetric: " + err.Error())
		}
		// Increment delta if current value exist
		if !errors.Is(err, storage.ErrNotFound) {
			// Check if current metric in not 'counter' type
			if m.MType != domain.MetricTypeCounter {
				return nil, ErrMTypeMismatch
			}
			// Add delta to current value
			*metric.Delta += *m.Delta
		}
	}
	// Flush hash value. We always store metrics without hash
	metric.Hash = ""
	// Write metric to storage
	m, err := s.storage.WriteMetric(ctx, metric)
	// Write metric to file if storeInterval == 0s
	if s.config.StoreInterval == 0*time.Second {
		err1 := s.fStorage.WriteMetric(ctx, metric)
		if err1 != nil {
			err = errors.Wrap(err, err1.Error())
		}
	}
	return m, err
}

func (s *Service) GetMetric(ctx context.Context, metric *domain.Metric) (*domain.Metric, error) {
	curMetric, err := s.storage.GetMetric(ctx, metric.ID)
	if curMetric != nil && curMetric.MType != metric.MType {
		return curMetric, ErrMTypeMismatch
	}
	return curMetric, err
}

func (s *Service) GetAllMetrics(ctx context.Context) (*domain.Metrics, error) {
	return s.storage.GetAllMetrics(ctx)
}

func (s *Service) WriteMetrics(ctx context.Context, metrics *domain.Metrics) error {
	for _, metric := range metrics.Metrics {
		_, err := s.WriteMetric(ctx, metric)
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
		metrics, err := s.fStorage.GetAllMetrics(context.Background())
		if err != nil {
			log.Printf("error occured while reading metrics from file: %v", err)
			return
		}
		err = s.storage.WriteAllMetrics(context.Background(), metrics)
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
	metrics, err := s.storage.GetAllMetrics(context.Background())
	if err != nil {
		return errors.New("storage.GetAllMetrics: error occurred while reading all metrics from storage: " + err.Error())
	}
	// Write all metrics to file
	err = s.fStorage.WriteAllMetrics(context.Background(), metrics)
	if err != nil {
		return errors.New("fStorage.WriteAllMetrics: error occurred while dumping data to file: " + err.Error())
	}
	return nil
}

func (s *Service) Ping() error {
	return s.storage.PingRepo()
}

func (s *Service) Shutdown() error {
	if s.config.StoreFile != "" {
		err := s.WriteAllToFile()
		if err != nil {
			return errors.New("WriteAllToFile: error occurred while dumping data to file: " + err.Error())
		}
		log.Printf("Stored to file before shutdown")
	}
	return nil
}

var ErrMTypeMismatch = errors.New("metric type mismatch")
