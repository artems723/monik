package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"

	"github.com/artems723/monik/internal/server/domain"
)

type FileStorage struct {
	encoder  *json.Encoder
	file     *os.File
	fileName string
}

func NewFileStorage(fileName string) *FileStorage {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalf("Error creating new file storage: %v", err)
	}
	return &FileStorage{
		file:     file,
		fileName: fileName,
		encoder:  json.NewEncoder(file),
	}
}

func (s *FileStorage) Close() error {
	return s.file.Close()
}

func (s *FileStorage) WriteAllMetrics(ctx context.Context, metrics *domain.Metrics) error {
	return s.encoder.Encode(metrics)
}

func (s *FileStorage) GetAllMetrics(ctx context.Context) (*domain.Metrics, error) {
	// read our opened jsonFile as a byte array.
	byteValue, err := s.readLastLine()
	if err != nil {
		return nil, errors.New("error reading last line of the file: " + err.Error())
	}

	var metrics domain.Metrics
	err = json.Unmarshal(byteValue, &metrics)
	if err != nil {
		return nil, err
	}
	return &metrics, nil
}

func (s *FileStorage) GetMetric(ctx context.Context, metricName string) (*domain.Metric, error) {
	metrics, err := s.GetAllMetrics(ctx)
	if err != nil {
		return nil, err
	}
	for _, val := range metrics.Metrics {
		if val.ID == metricName {
			return val, nil
		}
	}
	return nil, ErrNotFound
}

func (s *FileStorage) WriteMetric(ctx context.Context, metric *domain.Metric) error {
	return s.WriteAllMetrics(ctx, &domain.Metrics{Metrics: []*domain.Metric{metric}})
}

func (s *FileStorage) readLastLine() ([]byte, error) {

	reader := bufio.NewReader(s.file)

	// calculate the size of the last line
	var lastLineSize int
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		lastLineSize = len(line)
	}
	// check if file is empty
	if lastLineSize == 0 {
		return nil, ErrEmptyFile
	}
	fileInfo, err := os.Stat(s.fileName)
	if err != nil {
		return nil, err
	}
	// +1 to compensate for the initial 0 byte of the line
	buffer := make([]byte, lastLineSize)
	// read file from certain offset
	offset := fileInfo.Size() - int64(lastLineSize+1)
	numRead, err := s.file.ReadAt(buffer, offset)
	if err != nil {
		return nil, err
	}
	buffer = buffer[:numRead]
	return buffer, nil
}

func (s *FileStorage) PingRepo() error {
	_, err := s.file.Stat()
	return err
}

var ErrEmptyFile = errors.New("file is empty")
