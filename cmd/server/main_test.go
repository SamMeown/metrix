package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
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
			storage := MemStorage{make(map[string]any)}
			req := httptest.NewRequest(tt.requestMethod, tt.requestPath, nil)
			recorder := httptest.NewRecorder()
			handler := metricsRouter(storage)

			handler.ServeHTTP(recorder, req)

			result := recorder.Result()
			io.Copy(io.Discard, result.Body)
			result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
		})
	}
}
