package middleware

import (
	"github.com/SamMeown/metrix/internal/logger"
	"net/http"
	"time"
)

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func Logging(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		url := r.RequestURI
		method := r.Method

		response := &responseData{}
		lResponseWriter := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   response,
		}

		next.ServeHTTP(&lResponseWriter, r)

		duration := time.Since(start)

		logger.Log.Infoln(
			"url", url,
			"method", method,
			"status", response.status,
			"duration", duration,
			"size", response.size,
		)
	}

	return http.HandlerFunc(fn)
}
