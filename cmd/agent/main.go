package main

import (
	"fmt"
	"github.com/SamMeown/metrix/internal/agent/config"
	"io"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

type gauge = float64
type counter = int64

var memStatsMetricsNames = [...]string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}

type metricsCollection struct {
	gaugeMetrics   map[string]gauge
	counterMetrics map[string]counter
}

type metricsGetter struct {
	getCount             int64
	memStatsMetricsNames []string
}

func (mg *metricsGetter) getMetrics() metricsCollection {
	mg.getCount++
	metrics := metricsCollection{
		gaugeMetrics:   mg.getMemstatMetrics(),
		counterMetrics: make(map[string]counter),
	}

	metrics.gaugeMetrics["RandomValue"] = float64(rand.Int())
	metrics.counterMetrics["PollCount"] = mg.getCount

	return metrics
}

func (mg metricsGetter) getMemstatMetrics() map[string]gauge {
	metrics := make(map[string]gauge)
	var memstats runtime.MemStats
	runtime.ReadMemStats(&memstats)

	for _, name := range mg.memStatsMetricsNames {
		memstatsValue := reflect.ValueOf(memstats)
		field := memstatsValue.FieldByName(name)

		if field.CanFloat() {
			metrics[name] = gauge(field.Float())
		} else if field.CanUint() {
			metrics[name] = gauge(field.Uint())
		} else {
			metrics[name] = gauge(field.Int())
		}
	}

	return metrics
}

func (mg *metricsGetter) resetCount() {
	mg.getCount = 0
}

type metricsClient struct {
	http.Client
	baseURL string
}

func (client *metricsClient) reportAllMetrics(metrics metricsCollection) {
	for name, value := range metrics.gaugeMetrics {
		err := client.reportMetrics(name, value)
		if err != nil {
			fmt.Println(err)
		}
	}

	for name, value := range metrics.counterMetrics {
		err := client.reportMetrics(name, value)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (client *metricsClient) reportMetrics(name string, value any) error {
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

func main() {
	agentConfig := config.Parse()

	getter := metricsGetter{memStatsMetricsNames: memStatsMetricsNames[:]}
	client := &metricsClient{baseURL: fmt.Sprintf("http://%s/update", agentConfig.ServerBaseAddress)}
	var metrics metricsCollection
	var secondsElapsed int64

	for {
		if secondsElapsed%int64(agentConfig.PollInterval) == 0 {
			fmt.Println("Refreshing metrics...")
			metrics = getter.getMetrics()
		}

		if secondsElapsed%int64(agentConfig.ReportInterval) == 0 {
			fmt.Println("Reporting metrics...")
			client.reportAllMetrics(metrics)
			getter.resetCount()
		}

		secondsElapsed++
		time.Sleep(1 * time.Second)
	}
}
