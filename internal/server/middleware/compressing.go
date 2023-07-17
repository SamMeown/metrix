package middleware

import (
	"compress/gzip"
	"io"
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

type gzipReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newGzipReader(r io.ReadCloser) (*gzipReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &gzipReader{
		r:  r,
		zr: zr,
	}, nil
}

func (g *gzipReader) Read(p []byte) (int, error) {
	return g.zr.Read(p)
}

func (g *gzipReader) Close() error {
	if err := g.r.Close(); err != nil {
		return err
	}

	return g.zr.Close()
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

		contentEncoding := req.Header.Get("Content-Encoding")
		gzipEncoded := strings.HasPrefix(contentEncoding, "gzip")
		if gzipEncoded {
			gr, err := newGzipReader(req.Body)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer gr.Close()

			req.Body = gr
		}

		next.ServeHTTP(res, req)
	}

	return http.HandlerFunc(fn)
}
