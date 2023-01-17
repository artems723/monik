package service

import (
	"bufio"
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
	file     *os.File
	fileName string
	encoder  *json.Encoder
	decoder  *json.Decoder
	repo     storage.Repository
}

func NewStore(fileName string, storage storage.Repository) (*Store, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &Store{
		file:     file,
		fileName: fileName,
		encoder:  json.NewEncoder(file),
		decoder:  json.NewDecoder(file),
		repo:     storage,
	}, nil
}

func (s *Store) Close() error {
	err := s.WriteMetrics()
	if err != nil {
		log.Printf("error occured while dumping data to file: %v", err)
		return err
	}
	log.Printf("Stored to file before shutdown")
	return nil
}

func (s *Store) Init(restore bool) {
	// Read metrics from file to storage
	if restore {
		metrics, err := s.ReadMetrics()
		if err != nil {
			log.Printf("error occured while reading metrics from file: %v", err)
			return
		}
		err = s.repo.WriteAllMetrics(metrics)
		if err != nil {
			log.Printf("error occured while writing metrics to storage: %v", err)
			return
		}
		log.Printf("The following metrics were loaded from file: %v", metrics)
	}
}

func (s *Store) Run(storeInterval time.Duration) {
	// infinite loop for dumping data to file
	storeIntervalTicker := time.NewTicker(storeInterval)
	for {
		select {
		case <-storeIntervalTicker.C:
			err := s.WriteMetrics()
			if err != nil {
				log.Printf("error occured while dumping data to file: %v", err)
				return
			}
			log.Printf("Stored to file")
		}
	}
}

func (s *Store) WriteMetrics() error {
	metrics, err := s.repo.GetAllMetrics()
	if err != nil {
		log.Printf("GetAllMetrics(), error: %v", err)
		return err
	}
	return s.encoder.Encode(metrics)
}

func (s *Store) ReadMetrics() (*domain.Metrics, error) {
	// read our opened jsonFile as a byte array.
	//byteValue, _ := io.ReadAll(s.file)

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

func (s *Store) readLastLine() ([]byte, error) {

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

var ErrEmptyFile = errors.New("file is empty")
