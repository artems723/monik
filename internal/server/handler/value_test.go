package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/artems723/monik/internal/server/config"
	"github.com/artems723/monik/internal/server/domain"
	"github.com/artems723/monik/internal/server/service"
	"github.com/artems723/monik/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestHandler_getValue(t *testing.T) {
	type fields struct {
		s          service.Service
		valGauge   float64
		valCounter int64
	}
	type want struct {
		contentType string
		metricValue string
		statusCode  int
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
		args      args
		fields    fields
		name      string
		urlParams urlParams
		want      want
	}{
		{
			name:      "test get gauge value",
			fields:    fields{s: *service.New(storage.NewMemStorage(), config.Config{StoreInterval: 1 * time.Second}), valGauge: 20.201},
			args:      args{httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/{metricType}/{metricName}", nil)},
			want:      want{contentType: "text/plain; charset=utf-8", statusCode: 200, metricValue: "20.201"},
			urlParams: urlParams{"gauge", "Alloc"},
		},
		{
			name:      "test get counter value",
			fields:    fields{s: *service.New(storage.NewMemStorage(), config.Config{StoreInterval: 1 * time.Second}), valCounter: 20},
			args:      args{httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/{metricType}/{metricName}", nil)},
			want:      want{contentType: "text/plain; charset=utf-8", statusCode: 200, metricValue: "20"},
			urlParams: urlParams{"counter", "PollCount"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				s: &tt.fields.s,
			}
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("metricType", tt.urlParams.metricType)
			rctx.URLParams.Add("metricName", tt.urlParams.metricName)

			tt.args.r = tt.args.r.WithContext(context.WithValue(tt.args.r.Context(), chi.RouteCtxKey, rctx))

			// add metric to storage
			var metric *domain.Metric
			switch domain.MetricType(tt.urlParams.metricType) {
			case domain.MetricTypeGauge:
				metric = domain.NewGaugeMetric(tt.urlParams.metricName, tt.fields.valGauge)
			case domain.MetricTypeCounter:
				metric = domain.NewCounterMetric(tt.urlParams.metricName, tt.fields.valCounter)
			}

			tt.fields.s.WriteMetric(tt.args.r.Context(), metric)

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
		metric      *domain.Metric
		statusCode  int
	}
	type args struct {
		contentType string
		metric      *domain.Metric
		r           *http.Request
		w           http.ResponseWriter
	}
	tests := []struct {
		args   args
		fields fields
		name   string
		want   want
	}{
		{
			name:   "test success path",
			fields: fields{s: *service.New(storage.NewMemStorage(), config.Config{StoreInterval: 1 * time.Second})},
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
				s: &tt.fields.s,
			}

			// add metric to storage
			tt.fields.s.WriteMetric(tt.args.r.Context(), tt.args.metric)

			// Set content-type
			tt.args.r.Header.Set("Content-Type", tt.args.contentType)
			// Run handler
			h.getValueJSON(tt.args.w, tt.args.r)
			// Get response
			response := tt.args.w.(*httptest.ResponseRecorder).Result()
			defer response.Body.Close()
			// Get JSON response as metric struct
			var b domain.Metric
			json.NewDecoder(response.Body).Decode(&b)

			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.statusCode, response.StatusCode)
			assert.Equal(t, *tt.want.metric, b)
		})
	}
}
