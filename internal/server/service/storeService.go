package service

import (
	"encoding/json"
	"errors"
	"github.com/artems723/monik/internal/server/domain"
	"github.com/artems723/monik/internal/server/storage"
	"io"
	"log"
	"os"
	"time"
)

type Store struct {
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
	repo    storage.Repository
}

func NewStore(fileName string, storage storage.Repository) (*Store, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &Store{
		file:    file,
		encoder: json.NewEncoder(file),
		decoder: json.NewDecoder(file),
		repo:    storage,
	}, nil
}

func (s *Store) Close() error {
	return s.file.Close()
}

func (s *Store) Init(restore bool) {
	// Read metrics from file to storage
	if restore {
		metrics, err := s.ReadMetrics()
		if err != nil {
			log.Printf("error occured while reading metrics from file: %v", err)
			return
		}
		err = s.WriteMetrics(metrics)
		if err != nil {
			log.Printf("error occured while writing metrics to storage: %v", err)
			return
		}
	}
}

func (s *Store) Run(storeInterval time.Duration) {
	// infinite loop for dumping data to file
	storeIntervalTicker := time.NewTicker(storeInterval)
	for {
		select {
		case <-storeIntervalTicker.C:
			metrics, err := s.repo.GetAllMetrics()
			if err != nil {
				log.Printf("GetAllMetrics(), error: %v", err)
			}
			err = s.WriteMetrics(metrics)
			if err != nil {
				log.Printf("error occured while dumping data to file: %v", err)
				return
			}
			log.Printf("stored")
		}
	}
}

func (s *Store) WriteMetrics(metrics *domain.Metrics) error {
	// TODO: check
	return s.encoder.Encode(&metrics)
}

func (s *Store) ReadMetrics() (*domain.Metrics, error) {
	// read our opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(s.file)

	var metrics domain.Metrics

	err := json.Unmarshal(byteValue, &metrics)
	if err != nil {
		return nil, err
	}
	return &metrics, nil
}

var ErrEmptyFile = errors.New("no file or empty file")
