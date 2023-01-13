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

func TestHandler_updateMetric(t *testing.T) {
	type fields struct {
		s storage.Repository
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
		name      string
		fields    fields
		args      args
		want      want
		urlParams urlParams
	}{
		{
			name:      "test 200 code",
			fields:    fields{storage.NewMemStorage()},
			args:      args{httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/{metricType}/{metricName}/{metricValue}", nil)},
			want:      want{"", 200},
			urlParams: urlParams{"counter", "name", "2"},
		},
		{
			name:      "test 400 code",
			fields:    fields{storage.NewMemStorage()},
			args:      args{httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/{metricType}/{metricName}/{metricValue}", nil)},
			want:      want{"text/plain; charset=utf-8", 400},
			urlParams: urlParams{"counter", "name", ""},
		},
		{
			name:      "test 200 code",
			fields:    fields{storage.NewMemStorage()},
			args:      args{httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/{metricType}/{metricName}/{metricValue}", nil)},
			want:      want{"", 200},
			urlParams: urlParams{"gauge", "name", "2"},
		},
		{
			name:      "test 400 code",
			fields:    fields{storage.NewMemStorage()},
			args:      args{httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/{metricType}/{metricName}/{metricValue}", nil)},
			want:      want{"text/plain; charset=utf-8", 400},
			urlParams: urlParams{"gauge", "name", ""},
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
