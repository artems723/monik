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
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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
			name:   "test main page",
			args:   args{httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/}", nil)},
			want:   want{"text/html; charset=utf-8", 200, "<!DOCTYPE html>\r\n<html lang=\"en\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <title>All metrics</title>\r\n</head>\r\n<body>\r\n<h1>All metrics</h1>\r\nName: Alloc, Type: gauge, Value: 20.200000<br>\r\n</body>\r\n</html>", "Alloc", 20.20},
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

			//wd, _ := os.Getwd()
			//fmt.Println(wd)
			//for !strings.HasSuffix(wd, "monik") {
			//	wd = filepath.Dir(wd)
			//}
			//fmt.Println(wd)

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
