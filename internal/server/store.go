package server

import (
	"encoding/json"
	"github.com/artems723/monik/internal/server/domain"
	"log"
	"os"
	"time"
)

type Store struct {
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
}

func NewStore(fileName string) (*Store, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &Store{
		file:    file,
		encoder: json.NewEncoder(file),
		decoder: json.NewDecoder(file),
	}, nil
}

func (s *Store) Close() error {
	return s.file.Close()
}

func (s *Store) Run(storeInterval time.Duration) error {
	// infinite loop for polling counters and sending it to server
	storeIntervalTicker := time.NewTicker(storeInterval)
	for {
		select {
		case <-storeIntervalTicker.C:
			log.Printf("store")
		}
	}
}

func (s *Store) WriteMetrics(metrics map[string]*domain.Metrics) error {
	// TODO: check
	return s.encoder.Encode(&metrics)
}

func (s *Store) ReadMetrics() ([]*domain.Metrics, error) {
	metrics := make([]*domain.Metrics, 30)
	// read open bracket
	_, err := s.decoder.Token()
	if err != nil {
		log.Fatal(err)
	}
	// while the array contains values
	for s.decoder.More() {
		var m domain.Metrics
		// decode an array value (Metrics)
		err := s.decoder.Decode(&m)
		// add metric to slice
		metrics = append(metrics, &m)
		if err != nil {
			log.Fatal(err)
		}
	}
	// read closing bracket
	_, err = s.decoder.Token()
	if err != nil {
		log.Fatal(err)
	}
	return metrics, nil
}
