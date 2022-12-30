package handler

import (
	"github.com/artems723/monik/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_mainPage(t *testing.T) {
	type fields struct {
		s  storage.Repository
		id string
	}
	type want struct {
		contentType string
		statusCode  int
		text        string
		metricName  string
		metricValue string
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
			name:   "test get value",
			fields: fields{s: storage.NewMemStorage()},
			args:   args{httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/}", nil)},
			want:   want{"text/plain; charset=utf-8", 200, "Alloc=\"20\"\n", "Alloc", "20"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				s: tt.fields.s,
			}

			// add metric to storage
			tt.fields.s.WriteMetric(tt.fields.id, tt.want.metricName, tt.want.metricValue)

			// change remote address
			tt.args.r.RemoteAddr = tt.fields.id

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
