package agent

import (
	"net/http"
	"testing"

	"github.com/SamMeown/metrix/internal/agent/client"
	"github.com/stretchr/testify/assert"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(transport RoundTripFunc) http.Client {
	return http.Client{
		Transport: transport,
	}
}

func TestMetricsClientRequest(t *testing.T) {
	tests := []struct {
		name         string
		metricsName  string
		metricsValue any
		wantedPath   string
	}{
		{
			name:         "Test send gauge metrics",
			metricsName:  "a",
			metricsValue: float64(1234.1234),
			wantedPath:   "/update/gauge/a/1234.1234",
		},
		{
			name:         "Test send counter metrics",
			metricsName:  "b",
			metricsValue: int64(1234),
			wantedPath:   "/update/counter/b/1234",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testClient := NewTestClient(func(req *http.Request) *http.Response {
				assert.Equal(t, req.Method, http.MethodPost)
				assert.Equal(t, req.URL.Path, tt.wantedPath)
				return &http.Response{
					StatusCode: http.StatusOK,
				}
			})

			mClient := client.NewMetricsCustomClient("localhost:8080", testClient)

			mClient.ReportMetricsV1(tt.metricsName, tt.metricsValue)
		})
	}
}
