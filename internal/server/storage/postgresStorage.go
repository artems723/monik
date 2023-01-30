package storage

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"github.com/artems723/monik/internal/server/domain"
	"github.com/jmoiron/sqlx"
	"log"
	"path/filepath"
	"time"
)

type PostgresStorage struct {
	db *sqlx.DB
}

//go:embed sql/*
var SQL embed.FS

func NewPostgresStorage(databaseDSN string) (*PostgresStorage, error) {
	db, err := sqlx.Connect("pgx", databaseDSN)
	if err != nil {
		return nil, err
	}

	// Create table if not exists
	path := filepath.Join("sql", "metrics_table_up.sql")
	file, err := SQL.ReadFile(path)
	if err != nil {
		return nil, err
	}
	schema := string(file)
	db.MustExec(schema)

	log.Printf("succesfully connected to postgres")
	return &PostgresStorage{db}, nil
}

func (p *PostgresStorage) GetMetric(ctx context.Context, metricName string) (*domain.Metric, error) {
	var m domain.Metric
	err := p.db.Get(&m, "SELECT name,type,delta,value FROM metrics WHERE name=$1", metricName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &m, nil
}

func (p *PostgresStorage) WriteMetric(ctx context.Context, metric *domain.Metric) error {
	tx := p.db.MustBegin()
	_, err := tx.NamedExec("INSERT INTO metrics (name, type, delta, value) VALUES (:name, :type, :delta, :value) ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name, type = EXCLUDED.type, delta = EXCLUDED.delta, value = EXCLUDED.value", &metric)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresStorage) GetAllMetrics(ctx context.Context) (*domain.Metrics, error) {
	var m []*domain.Metric
	err := p.db.Select(&m, "SELECT name,type,delta,value FROM metrics")
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &domain.Metrics{Metrics: m}, nil
}

func (p *PostgresStorage) WriteAllMetrics(ctx context.Context, metrics *domain.Metrics) error {
	tx := p.db.MustBegin()
	_, err := tx.NamedExec("INSERT INTO metrics (name, type, delta, value) VALUES (:name, :type, :delta, :value) ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name, type = EXCLUDED.type, delta = EXCLUDED.delta, value = EXCLUDED.value", &metrics.Metrics)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresStorage) PingRepo() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return p.db.PingContext(ctx)
}
