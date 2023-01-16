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

func (s *Store) WriteMetrics(metrics domain.Metric) error {
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
