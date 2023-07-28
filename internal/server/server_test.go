package server

import (
	"bytes"
	"context"
	"github.com/SamMeown/metrix/internal/storage/mock"
	"github.com/golang/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SamMeown/metrix/internal/server/config"
	"github.com/SamMeown/metrix/internal/server/saver"
	"github.com/SamMeown/metrix/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestHandleUpdate(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
	}
	tests := []struct {
		name          string
		requestMethod string
		requestPath   string
		want          want
	}{
		{
			name:          "test no metrics type",
			requestMethod: http.MethodPost,
			requestPath:   "/update/",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:          "test wrong metrics type",
			requestMethod: http.MethodPost,
			requestPath:   "/update/gaga",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:          "test no metrics name",
			requestMethod: http.MethodPost,
			requestPath:   "/update/gauge",
			want: want{
				statusCode:  http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:          "test no metrics value",
			requestMethod: http.MethodPost,
			requestPath:   "/update/gauge/a",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:          "test wrong metrics value",
			requestMethod: http.MethodPost,
			requestPath:   "/update/gauge/a/notanumber",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:          "test right gauge request",
			requestMethod: http.MethodPost,
			requestPath:   "/update/gauge/a/11.35",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:          "test right counter request",
			requestMethod: http.MethodPost,
			requestPath:   "/update/counter/a/11",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:          "test counter request with wrong (float) value",
			requestMethod: http.MethodPost,
			requestPath:   "/update/counter/a/11.35",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			testConfig := config.Config{
				StoreInterval: 999999,
				Restore:       false,
			}
			nullSaver := &saver.MetricsStorageSaver{}
			mStorage := storage.New()

			req := httptest.NewRequest(tt.requestMethod, tt.requestPath, nil)
			recorder := httptest.NewRecorder()
			handler := metricsRouter(context.Background(), testConfig, mStorage, nullSaver, nil)

			handler.ServeHTTP(recorder, req)

			result := recorder.Result()
			io.Copy(io.Discard, result.Body)
			result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
		})
	}
}

func TestHandleValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	mStorage := mock.NewMockMetricsStorage(ctrl)

	gaugeValue := float64(42)
	counterValue := int64(1)
	mStorage.EXPECT().GetGauge(gomock.Any(), "a").Return(&gaugeValue, nil)
	mStorage.EXPECT().GetCounter(gomock.Any(), "b").Return(&counterValue, nil)

	nullSaver := &saver.MetricsStorageSaver{}
	testConfig := config.Config{
		StoreInterval: 999999,
		Restore:       false,
	}

	type want struct {
		statusCode  int
		contentType string
		body        string
	}
	tests := []struct {
		name          string
		requestMethod string
		requestPath   string
		want          want
	}{
		{
			name:          "test right gauge request",
			requestMethod: http.MethodGet,
			requestPath:   "/value/gauge/a",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
				body:        "42",
			},
		},
		{
			name:          "test right counter request",
			requestMethod: http.MethodGet,
			requestPath:   "/value/counter/b",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
				body:        "1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.requestMethod, tt.requestPath, nil)
			recorder := httptest.NewRecorder()
			handler := metricsRouter(context.Background(), testConfig, mStorage, nullSaver, nil)

			handler.ServeHTTP(recorder, req)

			body := &bytes.Buffer{}
			result := recorder.Result()
			io.Copy(body, result.Body)
			result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			if len(tt.want.body) > 0 {
				assert.Equal(t, tt.want.body, strings.TrimSuffix(body.String(), "\n"))
			}
		})
	}
}
