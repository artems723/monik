package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/artems723/monik/internal/server/domain"
	"github.com/artems723/monik/internal/server/service"
	"github.com/artems723/monik/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_getValue(t *testing.T) {
	type fields struct {
		s          service.Service
		id         string
		valGauge   float64
		valCounter int64
	}
	type want struct {
		contentType string
		statusCode  int
		metricValue string
	}
	type urlParams struct {
		metricType string
		metricName string
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      want
		urlParams urlParams
	}{
		{
			name:      "test get gauge value",
			fields:    fields{s: service.New(storage.NewMemStorage()), valGauge: 20.201},
			args:      args{httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/{metricType}/{metricName}", nil)},
			want:      want{"text/plain; charset=utf-8", 200, "20.201"},
			urlParams: urlParams{"gauge", "Alloc"},
		},
		{
			name:      "test get counter value",
			fields:    fields{s: service.New(storage.NewMemStorage()), valCounter: 20},
			args:      args{httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/{metricType}/{metricName}", nil)},
			want:      want{"text/plain; charset=utf-8", 200, "20"},
			urlParams: urlParams{"counter", "PollCount"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				s: tt.fields.s,
			}
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("metricType", tt.urlParams.metricType)
			rctx.URLParams.Add("metricName", tt.urlParams.metricName)

			tt.args.r = tt.args.r.WithContext(context.WithValue(tt.args.r.Context(), chi.RouteCtxKey, rctx))

			// add metric to storage
			var metric *domain.Metrics
			switch domain.MetricType(tt.urlParams.metricType) {
			case domain.MetricTypeGauge:
				metric = domain.NewGaugeMetric(tt.urlParams.metricName, tt.fields.valGauge)
			case domain.MetricTypeCounter:
				metric = domain.NewCounterMetric(tt.urlParams.metricName, tt.fields.valCounter)
			}

			tt.fields.s.WriteMetric(tt.fields.id, metric)

			// change remote address
			tt.args.r.RemoteAddr = tt.fields.id

			// handler call
			h.getValue(tt.args.w, tt.args.r)
			response := tt.args.w.(*httptest.ResponseRecorder).Result()
			defer response.Body.Close()
			b, err := io.ReadAll(response.Body)
			if err != nil {
				log.Fatalln(err)
			}
			assert.Equal(t, tt.want.metricValue, string(b))
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.statusCode, response.StatusCode)
		})
	}
}

func TestHandler_getValueJSON(t *testing.T) {
	type fields struct {
		s service.Service
	}
	type want struct {
		contentType string
		statusCode  int
		metric      *domain.Metrics
	}
	type args struct {
		w           http.ResponseWriter
		r           *http.Request
		contentType string
		metric      *domain.Metrics
		id          string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
		args   args
	}{
		{
			name:   "test success path",
			fields: fields{s: service.New(storage.NewMemStorage())},
			want: want{
				contentType: "application/json",
				statusCode:  200,
				metric:      domain.NewCounterMetric("PollCount", 6),
			},
			args: args{
				w:           httptest.NewRecorder(),
				r:           httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("{\"id\":\"PollCount\",\"type\":\"counter\",\"delta\":6}"))),
				contentType: "application/json",
				metric:      domain.NewCounterMetric("PollCount", 6),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				s: tt.fields.s,
			}

			// add metric to storage
			tt.fields.s.WriteMetric(tt.args.id, tt.args.metric)

			// Set content-type
			tt.args.r.Header.Set("Content-Type", tt.args.contentType)
			// change remote address
			tt.args.r.RemoteAddr = tt.args.id
			// Run handler
			h.getValueJSON(tt.args.w, tt.args.r)
			// Get response
			response := tt.args.w.(*httptest.ResponseRecorder).Result()
			defer response.Body.Close()
			// Get JSON response as metric struct
			var b domain.Metrics
			json.NewDecoder(response.Body).Decode(&b)

			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.statusCode, response.StatusCode)
			assert.Equal(t, *tt.want.metric, b)
		})
	}
}
