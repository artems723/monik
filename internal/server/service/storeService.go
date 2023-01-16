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
	storage storage.Repository
}

func NewStoreService(fileName string, storage storage.Repository) (*Store, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &Store{
		file:    file,
		encoder: json.NewEncoder(file),
		decoder: json.NewDecoder(file),
		storage: storage,
	}, nil
}

func (s *Store) Close() error {
	return s.file.Close()
}

func (s *Store) Run(storeInterval time.Duration, storeFile string, restore bool) error {

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
	if err == io.EOF {
		return nil, ErrEmptyFile
	}
	if err != nil {
		return nil, errors.New("error reading file")
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

var ErrEmptyFile = errors.New("no file or empty file")
