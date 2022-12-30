package storage

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestMemStorage_GetMetric(t *testing.T) {
	type fields struct {
		storage map[string]map[string]string
	}
	type args struct {
		agentID    string
		metricName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  bool
	}{
		{
			name:   "test read",
			fields: fields{storage: NewMemStorage().storage},
			args:   args{agentID: "127.0.0.1", metricName: "testMetric"},
			want:   "2",
			want1:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				storage: tt.fields.storage,
			}
			m.storage[tt.args.agentID] = make(map[string]string)
			m.storage[tt.args.agentID][tt.args.metricName] = "2"
			got, got1 := m.GetMetric(tt.args.agentID, tt.args.metricName)
			if got != tt.want {
				t.Errorf("GetMetric() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetMetric() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMemStorage_WriteMetric(t *testing.T) {
	type fields struct {
		storage map[string]map[string]string
	}
	type args struct {
		agentID     string
		metricName  string
		metricValue string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "test write",
			fields: fields{storage: NewMemStorage().storage},
			args:   args{agentID: "127.0.0.1", metricName: "testMetric", metricValue: "2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				storage: tt.fields.storage,
			}
			m.WriteMetric(tt.args.agentID, tt.args.metricName, tt.args.metricValue)
			assert.Equal(t, m.storage[tt.args.agentID][tt.args.metricName], tt.args.metricValue)
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
			want: &MemStorage{storage: make(map[string]map[string]string)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMemStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMemStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}
