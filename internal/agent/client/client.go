package client

import (
	"fmt"
	"github.com/SamMeown/metrix/internal/storage"
	"io"
	"net/http"
	"strconv"
)

type gauge = float64
type counter = int64

type MetricsClient struct {
	http.Client
	baseURL string
}

func NewMetricsClient(baseURL string) *MetricsClient {
	return &MetricsClient{
		baseURL: fmt.Sprintf("http://%s/update", baseURL),
	}
}

func NewMetricsCustomClient(baseURL string, client http.Client) *MetricsClient {
	return &MetricsClient{
		Client:  client,
		baseURL: fmt.Sprintf("http://%s/update", baseURL),
	}
}

func (client *MetricsClient) ReportAllMetrics(metricsCollection storage.MetricsStorageGetter) {
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

func (client *MetricsClient) ReportMetrics(name string, value any) error {
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
