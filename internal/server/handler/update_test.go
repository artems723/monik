package handler

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestHandler_updateMetric(t *testing.T) {
	type fields struct {
		s service.Service
	}
	type want struct {
		contentType string
		statusCode  int
	}
	type urlParams struct {
		metricType  string
		metricName  string
		metricValue string
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
			name:      "test 200 code",
			fields:    fields{*service.New(storage.NewMemStorage(), config.Config{StoreInterval: 1 * time.Second})},
			args:      args{httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/{metricType}/{metricName}/{metricValue}", nil)},
			want:      want{"", 200},
			urlParams: urlParams{"counter", "name", "2"},
		},
		{
			name:      "test 400 code",
			fields:    fields{*service.New(storage.NewMemStorage(), config.Config{StoreInterval: 1 * time.Second})},
			args:      args{httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/{metricType}/{metricName}/{metricValue}", nil)},
			want:      want{"text/plain; charset=utf-8", 400},
			urlParams: urlParams{"counter", "name", ""},
		},
		{
			name:      "test 200 code",
			fields:    fields{*service.New(storage.NewMemStorage(), config.Config{StoreInterval: 1 * time.Second})},
			args:      args{httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/{metricType}/{metricName}/{metricValue}", nil)},
			want:      want{"", 200},
			urlParams: urlParams{"gauge", "name", "2"},
		},
		{
			name:      "test 400 code",
			fields:    fields{*service.New(storage.NewMemStorage(), config.Config{StoreInterval: 1 * time.Second})},
			args:      args{httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/{metricType}/{metricName}/{metricValue}", nil)},
			want:      want{"text/plain; charset=utf-8", 400},
			urlParams: urlParams{"gauge", "name", ""},
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
			rctx.URLParams.Add("metricValue", tt.urlParams.metricValue)

			tt.args.r = tt.args.r.WithContext(context.WithValue(tt.args.r.Context(), chi.RouteCtxKey, rctx))

			h.updateMetric(tt.args.w, tt.args.r)
			response := tt.args.w.(*httptest.ResponseRecorder).Result()
			defer response.Body.Close()
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.statusCode, response.StatusCode)
		})
	}
}

func TestHandler_updateMetricJSON(t *testing.T) {
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
			// Set content-type
			tt.args.r.Header.Set("Content-Type", tt.args.contentType)
			// Run handler
			h.updateMetricJSON(tt.args.w, tt.args.r)
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
