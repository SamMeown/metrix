package middleware

import (
	"bytes"
	"github.com/SamMeown/metrix/internal/crypto/signer"
	"github.com/SamMeown/metrix/internal/logger"
	"io"
	"net/http"
)

type deferredWriter struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (d *deferredWriter) WriteHeader(status int) {
	d.status = status
}

func (d *deferredWriter) Write(body []byte) (int, error) {
	return d.body.Write(body)
}

func Signing(signer *signer.Signer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(res http.ResponseWriter, req *http.Request) {
			dr := &deferredWriter{
				ResponseWriter: res,
			}

			next.ServeHTTP(dr, req)

			if dr.body.Len() > 0 {
				signature := signer.GetSignature(dr.body.Bytes())
				res.Header().Add("HashSHA256", signature)
			}

			if dr.status != 0 {
				res.WriteHeader(dr.status)
			}

			_, err := io.Copy(res, &dr.body)
			if err != nil {
				http.Error(res, "Failed to write response body", http.StatusInternalServerError)
				return
			}
		}

		return http.HandlerFunc(fn)
	}
}

func SignValidating(signer *signer.Signer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(res http.ResponseWriter, req *http.Request) {
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				http.Error(res, "Failed to read body", http.StatusInternalServerError)
				return
			}

			signature := req.Header.Get("HashSHA256")
			if signature == "" && len(bodyBytes) > 0 {
				logger.Log.Debugln("Content signature not found")
				//http.Error(res, "Content signature not found", http.StatusBadRequest)
				//return
			}
			if signature != "" && !signer.ValidateSignature(signature, bodyBytes) {
				http.Error(res, "Content signature is not valid", http.StatusBadRequest)
				return
			}

			req.Body.Close()
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			next.ServeHTTP(res, req)
		}

		return http.HandlerFunc(fn)
	}
}
