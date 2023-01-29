package storage

import (
	"context"
	"database/sql"
	"github.com/artems723/monik/internal/server/domain"
	"log"
	"time"
)

type PostgresStorage struct {
	*sql.DB
}

func NewPostgresStorage(databaseDSN string) (*PostgresStorage, error) {
	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		return nil, err
	}
	// Check connection
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	log.Printf("succesfully connected to postgres")
	return &PostgresStorage{db}, nil
}

func (p *PostgresStorage) GetMetric(ctx context.Context, metricName string) (*domain.Metric, error) {
	// TODO
	return nil, nil
}

func (p *PostgresStorage) WriteMetric(ctx context.Context, metric *domain.Metric) error {
	// TODO
	return nil
}

func (p *PostgresStorage) GetAllMetrics(ctx context.Context) (*domain.Metrics, error) {
	// TODO
	return nil, nil
}

func (p *PostgresStorage) WriteAllMetrics(ctx context.Context, metrics *domain.Metrics) error {
	// TODO
	return nil
}

func (p *PostgresStorage) PingRepo() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return p.DB.PingContext(ctx)
}
