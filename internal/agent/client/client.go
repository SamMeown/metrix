package client

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SamMeown/metrix/internal/backoff"
	"io"
	"net"
	"net/http"
	"strconv"

	"github.com/SamMeown/metrix/internal/logger"
	"github.com/SamMeown/metrix/internal/models"
	"github.com/SamMeown/metrix/internal/storage"
)

type gauge = float64
type counter = int64

type MetricsClient struct {
	http.Client
	baseURL string
}

func NewMetricsClient(baseURL string) *MetricsClient {
	return &MetricsClient{
		baseURL: fmt.Sprintf("http://%s/updates", baseURL),
	}
}

func NewMetricsCustomClient(baseURL string, client http.Client) *MetricsClient {
	return &MetricsClient{
		Client:  client,
		baseURL: fmt.Sprintf("http://%s/update", baseURL),
	}
}

func NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	buf := bytes.Buffer{}
	gz, err := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(gz, body)
	if err != nil {
		return nil, err
	}
	gz.Close()

	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return req, err
	}

	req.Header.Set("Content-Encoding", "gzip")

	return req, nil
}

func (client *MetricsClient) ReportAllMetrics(metricsCollection storage.MetricsStorageGetter) {
	allMetrics, _ := metricsCollection.GetAll()
	metrics := make([]models.Metrics, 0)
	for name, value := range allMetrics.Gauges {
		reqMetrics, err := metricsToRequestMetrics(name, value)
		if err != nil {
			logger.Log.Errorln(err)
			return
		}
		metrics = append(metrics, reqMetrics)
	}
	for name, value := range allMetrics.Counters {
		reqMetrics, err := metricsToRequestMetrics(name, value)
		if err != nil {
			logger.Log.Errorln(err)
			return
		}
		metrics = append(metrics, reqMetrics)
	}

	logger.Log.Debugf("Reporting metrics: %+v", metrics)

	body, err := json.Marshal(metrics)
	if err != nil {
		logger.Log.Errorln(err)
		return
	}

	respCode, respBody, err := client.sendRequest(body)
	if err != nil {
		logger.Log.Errorln(err)
		return
	}

	logger.Log.Debugf("Status code: %d\nReport response: %s\n", respCode, respBody)
}

func metricsToRequestMetrics(name string, value any) (models.Metrics, error) {
	var metrics = models.Metrics{ID: name}

	switch typedValue := value.(type) {
	case gauge:
		metrics.MType = "gauge"
		metrics.Value = &typedValue
	case counter:
		metrics.MType = "counter"
		metrics.Delta = &typedValue
	default:
		return models.Metrics{}, errors.New("wrong metrics value type")
	}

	return metrics, nil
}

func (client *MetricsClient) sendRequest(requestBody []byte) (code int, body []byte, err error) {
	bOff := backoff.NewBackoff([]int{1, 3, 5})
	err = bOff.Retry(func() (e error) {
		code, body, e = client._sendRequest(requestBody)
		return
	})

	return
}

func (client *MetricsClient) _sendRequest(requestBody []byte) (code int, body []byte, err error) {
	req, err := NewRequest(http.MethodPost, client.baseURL, bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		var netErr *net.OpError
		if errors.As(err, &netErr) {
			err = backoff.NewRetryableError(err)
		}
		return
	}

	code = response.StatusCode

	defer response.Body.Close()
	body, err = io.ReadAll(response.Body)

	return
}

func (client *MetricsClient) ReportMetrics(name string, value any) error {
	metrics, err := metricsToRequestMetrics(name, value)
	if err != nil {
		panic(err)
	}

	logger.Log.Debugf("Reporting metrics: %+v", metrics)

	body, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	respCode, respBody, err := client.sendRequest(body)
	if err != nil {
		return err
	}

	logger.Log.Debugf("Status code: %d\nReport response: %s\n", respCode, respBody)

	return nil
}

func (client *MetricsClient) ReportAllMetricsV1(metricsCollection storage.MetricsStorageGetter) {
	allMetrics, _ := metricsCollection.GetAll()
	for name, value := range allMetrics.Gauges {
		err := client.ReportMetrics(name, value)
		if err != nil {
			fmt.Println(err)
		}
	}

	for name, value := range allMetrics.Counters {
		err := client.ReportMetrics(name, value)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (client *MetricsClient) ReportMetricsV1(name string, value any) error {
	var metricsType, valueString string
	switch typedValue := value.(type) {
	case gauge:
		metricsType = "gauge"
		valueString = strconv.FormatFloat(typedValue, 'f', -1, 64)
	case counter:
		metricsType = "counter"
		valueString = strconv.FormatInt(typedValue, 10)
	default:
		panic("Wrong metrics value type")
	}

	url := fmt.Sprintf("%s/%s/%s/%s", client.baseURL, metricsType, name, valueString)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "text/plain")

	response, err := client.Do(req)
	if err != nil {
		return err
	}
	fmt.Printf("Status code: %d\n", response.StatusCode)
	defer response.Body.Close()
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(respBody))

	return nil
}
