package client

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
)

func TestAgent_SendData(t *testing.T) {
	type fields struct {
		storage map[string]*Metric
	}
	type args struct {
		URL    string
		client HTTPClient
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "test send",
			fields: fields{storage: make(map[string]*Metric)},
			args:   args{URL: "", client: NewHTTPClient()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := &Agent{
				storage: tt.fields.storage,
				mu:      &sync.RWMutex{},
			}
			teardown := setup()
			defer teardown()
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			agent.storage["NumGC"] = NewGaugeMetric("NumGC", 222, "")
			agent.SendData(server.URL, tt.args.client)
		})
	}
}

func TestAgent_UpdateMetrics(t *testing.T) {
	type fields struct {
		metrics map[string]*Metric
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "test update metrics",
			fields: fields{metrics: make(map[string]*Metric)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := &Agent{
				storage: tt.fields.metrics,
				mu:      &sync.RWMutex{},
			}
			agent.UpdateMetrics()
			assert.Equal(t, *agent.storage["PollCount"].Delta, int64(1))
			_, ok := agent.storage["Alloc"]
			assert.Equal(t, true, ok)
		})
	}
}

func TestNewAgent(t *testing.T) {
	tests := []struct {
		name string
		want Agent
	}{
		{
			name: "test new agent",
			want: Agent{storage: make(map[string]*Metric)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAgent(""); !reflect.DeepEqual(got.storage, tt.want.storage) {
				t.Errorf("NewAgent() = %v, want %v", got, tt.want)
			}
		})
	}
}

var (
	mux    *http.ServeMux
	server *httptest.Server
)

func setup() func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	return func() {
		server.Close()
	}
}
