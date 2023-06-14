package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	MetricsTypeGauge   = "gauge"
	MetricsTypeCounter = "counter"
)

type MetricsStorage interface {
	SetGauge(name string, value float64) error
	SetCounter(name string, value int64) error
	Value(name string) (any, error)
}

type MemStorage struct {
	values map[string]any
}

func (m MemStorage) Value(name string) (any, error) {
	return m.values[name], nil
}

func (m MemStorage) SetGauge(name string, value float64) error {
	m.values[name] = value
	return nil
}

func (m MemStorage) SetCounter(name string, value int64) error {
	if _, ok := m.values[name]; !ok {
		m.values[name] = int64(0)
	}
	m.values[name] = m.values[name].(int64) + value

	return nil
}

var storage MetricsStorage = MemStorage{make(map[string]any)}

func handleUpdate(storage MetricsStorage) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		res.Header().Set("Content-Type", "text/plain; charset=utf-8")

		path := req.URL.Path
		fmt.Printf("Path: %s\n", path)

		var metricsType, metricsName string
		var metricsValue any

		components := strings.Split(strings.Trim(path, "/"), "/")[1:]
		//fmt.Printf("Components: %v", components)
		if len(components) < 1 || len(components) > 3 {
			http.Error(res, "Wrong number of data components", http.StatusBadRequest)
			return
		}

		metricsType = components[0]
		if metricsType != MetricsTypeGauge &&
			metricsType != MetricsTypeCounter {
			http.Error(res, "Wrong metrics type", http.StatusBadRequest)
			return
		}

		if len(components) < 2 {
			http.Error(res, "No metrics name", http.StatusNotFound)
			return
		}

		metricsName = components[1]

		if len(components) < 3 {
			http.Error(res, "No metrics value", http.StatusBadRequest)
			return
		}

		var convErr error
		if metricsType == MetricsTypeGauge {
			metricsValue, convErr = strconv.ParseFloat(components[2], 64)
		} else {
			metricsValue, convErr = strconv.ParseInt(components[2], 10, 64)
		}
		if convErr != nil {
			http.Error(res, "Can not parse metrics value", http.StatusBadRequest)
			return
		}

		//fmt.Printf("type: %s, name: %s, value: %v", metricsType, metricsName, metricsValue)
		if metricsType == MetricsTypeGauge {
			storage.SetGauge(metricsName, metricsValue.(float64))
		} else {
			storage.SetCounter(metricsName, metricsValue.(int64))
		}

		res.WriteHeader(http.StatusOK)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", handleUpdate(storage))

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
