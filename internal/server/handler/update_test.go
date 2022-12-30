package handler

import (
	"context"
	"github.com/artems723/monik/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_notImplemented(t *testing.T) {
	type fields struct {
		s storage.Repository
	}
	type want struct {
		contentType string
		statusCode  int
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name:   "test 501 code and content type",
			fields: fields{storage.NewMemStorage()},
			args:   args{httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/unknown", nil)},
			want:   want{"text/plain; charset=utf-8", 501},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				s: tt.fields.s,
			}
			h.notImplemented(tt.args.w, tt.args.r)
			response := tt.args.w.(*httptest.ResponseRecorder).Result()
			defer response.Body.Close()
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.statusCode, response.StatusCode)
		})
	}
}

func TestHandler_updateCounterMetric(t *testing.T) {
	type fields struct {
		s storage.Repository
	}
	type want struct {
		contentType string
		statusCode  int
	}
	type urlParams struct {
		metricName  string
		metricValue string
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
			name:      "test 200 code",
			fields:    fields{storage.NewMemStorage()},
			args:      args{httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/counter/{metricName}/{metricValue}", nil)},
			want:      want{"", 200},
			urlParams: urlParams{"name", "2"},
		},
		{
			name:      "test 400 code",
			fields:    fields{storage.NewMemStorage()},
			args:      args{httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/counter/{metricName}/{metricValue}", nil)},
			want:      want{"text/plain; charset=utf-8", 400},
			urlParams: urlParams{"name", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				s: tt.fields.s,
			}
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("metricName", tt.urlParams.metricName)
			rctx.URLParams.Add("metricValue", tt.urlParams.metricValue)

			tt.args.r = tt.args.r.WithContext(context.WithValue(tt.args.r.Context(), chi.RouteCtxKey, rctx))

			h.updateCounterMetric(tt.args.w, tt.args.r)
			response := tt.args.w.(*httptest.ResponseRecorder).Result()
			defer response.Body.Close()
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.statusCode, response.StatusCode)
		})
	}
}

func TestHandler_updateGaugeMetric(t *testing.T) {
	type fields struct {
		s storage.Repository
	}
	type want struct {
		contentType string
		statusCode  int
	}
	type urlParams struct {
		metricName  string
		metricValue string
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
			name:      "test 200 code",
			fields:    fields{storage.NewMemStorage()},
			args:      args{httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/gauge/{metricName}/{metricValue}", nil)},
			want:      want{"", 200},
			urlParams: urlParams{"name", "2"},
		},
		{
			name:      "test 400 code",
			fields:    fields{storage.NewMemStorage()},
			args:      args{httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/gauge/{metricName}/{metricValue}", nil)},
			want:      want{"text/plain; charset=utf-8", 400},
			urlParams: urlParams{"name", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				s: tt.fields.s,
			}
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("metricName", tt.urlParams.metricName)
			rctx.URLParams.Add("metricValue", tt.urlParams.metricValue)

			tt.args.r = tt.args.r.WithContext(context.WithValue(tt.args.r.Context(), chi.RouteCtxKey, rctx))

			h.updateGaugeMetric(tt.args.w, tt.args.r)
			response := tt.args.w.(*httptest.ResponseRecorder).Result()
			defer response.Body.Close()
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.statusCode, response.StatusCode)
		})
	}
}
