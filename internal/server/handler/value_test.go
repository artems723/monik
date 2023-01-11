package handler

import (
	"context"
	"github.com/artems723/monik/internal/server/domain"
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
		s          storage.Repository
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
			name:      "test get value",
			fields:    fields{s: storage.NewMemStorage(), valGauge: 20.201},
			args:      args{httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/{metricType}/{metricName}", nil)},
			want:      want{"text/plain; charset=utf-8", 200, "20.201"},
			urlParams: urlParams{"gauge", "Alloc"},
		},
		{
			name:      "test get value",
			fields:    fields{s: storage.NewMemStorage(), valCounter: 20},
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
			rctx.URLParams.Add("metricName", tt.urlParams.metricName)

			tt.args.r = tt.args.r.WithContext(context.WithValue(tt.args.r.Context(), chi.RouteCtxKey, rctx))

			// add metric to storage
			var metric domain.Metrics
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
