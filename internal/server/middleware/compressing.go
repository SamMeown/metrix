package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

var compressableTypes = []string{
	"application/json",
	"text/html",
}

type gzipWriter struct {
	w          http.ResponseWriter
	zw         *gzip.Writer
	compress   *bool
	statusCode *int
}

func newGzipWriter(w http.ResponseWriter) *gzipWriter {
	return &gzipWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (g *gzipWriter) shouldCompress() bool {
	if g.compress == nil {
		if g.statusCode == nil {
			panic("Unknown status code")
		}

		g.compress = new(bool)

		if *g.statusCode < 300 {
			contentType := g.w.Header().Get("Content-Type")
			for _, compressableType := range compressableTypes {
				if strings.Contains(contentType, compressableType) {
					*g.compress = true
					break
				}
			}
		}
	}

	return *g.compress
}

func (g *gzipWriter) Header() http.Header {
	return g.w.Header()
}

func (g *gzipWriter) WriteHeader(statusCode int) {
	g.statusCode = &statusCode
	if g.shouldCompress() {
		g.w.Header().Set("Content-Encoding", "gzip")
	}
	g.w.WriteHeader(statusCode)
}

func (g *gzipWriter) Write(p []byte) (int, error) {
	if g.statusCode == nil {
		g.WriteHeader(http.StatusOK)
	}
	if g.shouldCompress() {
		return g.zw.Write(p)
	} else {
		return g.w.Write(p)
	}
}

func (g *gzipWriter) Close() error {
	if g.shouldCompress() {
		return g.zw.Close()
	}

	return nil
}

func Compressing(next http.Handler) http.Handler {
	fn := func(res http.ResponseWriter, req *http.Request) {
		acceptEncoding := req.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			gzipRes := newGzipWriter(res)
			defer gzipRes.Close()
			res = gzipRes
		}

		next.ServeHTTP(res, req)
	}

	return http.HandlerFunc(fn)
}
