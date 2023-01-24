package handler

import (
	"github.com/artems723/monik/internal/server"
	"github.com/artems723/monik/internal/server/domain"
	"github.com/artems723/monik/internal/server/service"
	"github.com/artems723/monik/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_mainPage(t *testing.T) {
	type fields struct {
		s service.Service
	}
	type want struct {
		contentType string
		statusCode  int
		text        string
		metricName  string
		metricValue float64
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		fields fields
		name   string
		args   args
		want   want
	}{
		{
			fields: fields{s: *service.New(storage.NewMemStorage(), server.Config{StoreInterval: 1 * time.Second})},
			name:   "test get value",
			args:   args{httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/}", nil)},
			want:   want{"text/html; charset=utf-8", 200, "<!DOCTYPE html><html><body><h1>All metrics</h1></body>Name: Alloc, Type: gauge, Value: 20.200000<br></html>", "Alloc", 20.20},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				s: &tt.fields.s,
			}

			// add metric to storage
			metric := domain.NewGaugeMetric(tt.want.metricName, tt.want.metricValue)
			h.s.WriteMetric(metric)

			// handler call
			h.mainPage(tt.args.w, tt.args.r)
			response := tt.args.w.(*httptest.ResponseRecorder).Result()
			defer response.Body.Close()
			b, err := io.ReadAll(response.Body)
			if err != nil {
				log.Fatalln(err)
			}
			assert.Equal(t, tt.want.text, string(b))
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.statusCode, response.StatusCode)
		})
	}
}
