package storage

import (
	"github.com/artems723/monik/internal/server/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemStorage_GetMetric(t *testing.T) {
	type fields struct {
		storage map[string]*domain.Metric
	}
	type args struct {
		metricName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   domain.Metric
		want1  error
	}{
		{
			name:   "test read",
			fields: fields{storage: NewMemStorage().storage},
			args:   args{metricName: "testMetric"},
			want:   *domain.NewGaugeMetric("testMetric", 5.0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				storage: tt.fields.storage,
			}
			m.storage[tt.args.metricName] = domain.NewGaugeMetric("testMetric", 5.0)
			got, got1 := m.GetMetric(tt.args.metricName)

			assert.Equal(t, *got, tt.want)
			assert.Equal(t, got1, tt.want1)
		})
	}
}

func TestMemStorage_WriteMetric(t *testing.T) {
	type fields struct {
		storage map[string]*domain.Metric
	}
	type args struct {
		metric *domain.Metric
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "test write",
			fields: fields{storage: NewMemStorage().storage},
			args:   args{metric: domain.NewGaugeMetric("testMetric", 5.0)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				storage: tt.fields.storage,
			}
			m.WriteMetric(tt.args.metric)
			assert.Equal(t, m.storage[tt.args.metric.ID], tt.args.metric)
		})
	}
}

func TestNewMemStorage(t *testing.T) {
	tests := []struct {
		name string
		want *MemStorage
	}{
		{
			name: "test new storage",
			want: &MemStorage{storage: make(map[string]*domain.Metric)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMemStorage()
			assert.Equal(t, got, tt.want)
		})
	}
}
