package handler

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/artems723/monik/internal/server/config"
	"github.com/artems723/monik/internal/server/domain"
	"github.com/artems723/monik/internal/server/service"
	"github.com/artems723/monik/internal/server/storage"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Change working dir to root project dir. It is needed for locating template files.
	wd, _ := os.Getwd()
	for !strings.HasSuffix(wd, "monik") {
		wd = filepath.Dir(wd)
	}
	err := os.Chdir(wd)
	if err != nil {
		panic(err)
	}
}

func TestHandler_mainPage(t *testing.T) {
	type fields struct {
		s service.Service
	}
	type want struct {
		contentType string
		metricName  string
		metricValue float64
		statusCode  int
		text        string
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
			fields: fields{s: *service.New(storage.NewMemStorage(), config.Config{StoreInterval: 1 * time.Second})},
			name:   "test main page",
			args:   args{httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/}", nil)},
			want:   want{contentType: "text/html", statusCode: 200, text: "Name: Alloc, Type: gauge, Value: 20.200000", metricName: "Alloc", metricValue: 20.20},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New(&tt.fields.s, "", "")

			// add metric to storage
			metric := domain.NewGaugeMetric(tt.want.metricName, tt.want.metricValue)
			h.s.WriteMetric(tt.args.r.Context(), metric)

			// handler call
			h.mainPage(tt.args.w, tt.args.r)
			response := tt.args.w.(*httptest.ResponseRecorder).Result()
			defer response.Body.Close()
			b, err := io.ReadAll(response.Body)
			if err != nil {
				log.Fatalln(err)
			}
			assert.Contains(t, string(b), tt.want.text)
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.statusCode, response.StatusCode)
		})
	}
}
