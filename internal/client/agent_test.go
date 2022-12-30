package client

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestAgent_SendData(t *testing.T) {
	type fields struct {
		gaugeMetrics map[string]gauge
		pollCount    counter
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
			fields: fields{gaugeMetrics: make(map[string]gauge), pollCount: 2},
			args:   args{URL: "", client: NewHTTPClient()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := &Agent{
				gaugeMetrics: tt.fields.gaugeMetrics,
				pollCount:    tt.fields.pollCount,
			}
			teardown := setup()
			defer teardown()
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			agent.gaugeMetrics["NumGC"] = gauge(222)
			agent.SendData(server.URL, tt.args.client)
		})
	}
}

func TestAgent_UpdateMetrics(t *testing.T) {
	type fields struct {
		gaugeMetrics map[string]gauge
		pollCount    counter
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "test update metrics",
			fields: fields{gaugeMetrics: make(map[string]gauge), pollCount: 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := &Agent{
				gaugeMetrics: tt.fields.gaugeMetrics,
				pollCount:    tt.fields.pollCount,
			}
			agent.UpdateMetrics()
			assert.Equal(t, tt.fields.pollCount+1, agent.pollCount)
			_, ok := agent.gaugeMetrics["Alloc"]
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
			want: Agent{gaugeMetrics: make(map[string]gauge), pollCount: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAgent(); !reflect.DeepEqual(got, tt.want) {
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
