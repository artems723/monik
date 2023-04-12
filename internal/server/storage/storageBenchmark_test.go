package storage

import (
	"context"
	"github.com/artems723/monik/internal/server/domain"
	"os"
	"path/filepath"
	"testing"
)

func Benchmark_MemoryAndFileStorage(b *testing.B) {
	// setup memory storage
	m := NewMemStorage()
	m.s["testMetric"] = domain.NewGaugeMetric("testMetric", 5.0)
	// setup file storage
	path := filepath.Join(os.TempDir(), "devops-metrics-db.json")
	f := NewFileStorage(path)
	f.WriteMetric(context.Background(), domain.NewGaugeMetric("testMetric", 5.0))

	// GetAllMetrics
	b.Run("GetAllMetrics_MemStorage", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m.GetAllMetrics(context.Background())
		}
	})
	b.Run("GetAllMetrics_FileStorage", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			f.GetAllMetrics(context.Background())
		}
	})

	// GetMetric
	b.Run("GetMetric_MemStorage", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m.GetMetric(context.Background(), "testMetric")
		}
	})
	b.Run("GetMetric_FileStorage", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			f.GetMetric(context.Background(), "testMetric")
		}
	})

	// WriteMetric
	metric := domain.NewGaugeMetric("testMetric", 5.0)
	b.Run("WriteMetric_MemStorage", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m.WriteMetric(context.Background(), metric)
		}
	})
	b.Run("WriteMetric_FileStorage", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			f.WriteMetric(context.Background(), metric)
		}
	})

	// WriteAllMetrics
	metrics := domain.Metrics{Metrics: []*domain.Metric{domain.NewGaugeMetric("testMetric", 5.0)}}
	b.Run("WriteAllMetrics_MemStorage", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m.WriteAllMetrics(context.Background(), &metrics)
		}
	})
	b.Run("WriteAllMetrics_FileStorage", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			f.WriteAllMetrics(context.Background(), &metrics)
		}
	})
}
