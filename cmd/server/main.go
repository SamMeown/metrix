package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"strconv"
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
		res.Header().Set("Content-Type", "text/plain; charset=utf-8")

		path := req.URL.Path
		fmt.Printf("Path: %s\n", path)

		var metricsType = chi.URLParam(req, "metricsType")
		var metricsName = chi.URLParam(req, "metricsName")
		var metricsValueStr = chi.URLParam(req, "metricsValue")

		//fmt.Printf("metricsType: %s, metricsName: %s, metricsValue: %s\n", metricsType, metricsName, metricsValueStr)

		if metricsType == "" {
			http.Error(res, "Wrong number of data components", http.StatusBadRequest)
			return
		}

		if metricsType != MetricsTypeGauge &&
			metricsType != MetricsTypeCounter {
			http.Error(res, "Wrong metrics type", http.StatusBadRequest)
			return
		}

		if metricsName == "" {
			http.Error(res, "No metrics name", http.StatusNotFound)
			return
		}

		if metricsValueStr == "" {
			http.Error(res, "No metrics value", http.StatusBadRequest)
			return
		}

		var metricsValue any
		var convErr error
		if metricsType == MetricsTypeGauge {
			metricsValue, convErr = strconv.ParseFloat(metricsValueStr, 64)
		} else {
			metricsValue, convErr = strconv.ParseInt(metricsValueStr, 10, 64)
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

func handleValue(storage MetricsStorage) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/plain; charset=utf-8")

		path := req.URL.Path
		fmt.Printf("Path: %s\n", path)

		var metricsType = chi.URLParam(req, "metricsType")
		var metricsName = chi.URLParam(req, "metricsName")

		if metricsType != MetricsTypeGauge &&
			metricsType != MetricsTypeCounter {
			http.Error(res, "Wrong metrics type", http.StatusBadRequest)
			return
		}

		value, _ := storage.Value(metricsName)
		if value == nil {
			http.Error(res, "Metrics not found", http.StatusNotFound)
			return
		}

		var valueString string
		switch typedValue := value.(type) {
		case float64:
			if metricsType != MetricsTypeGauge {
				http.Error(res, "Metrics has another type", http.StatusNotFound)
				return
			}
			valueString = strconv.FormatFloat(typedValue, 'f', 2, 64)
		case int64:
			if metricsType != MetricsTypeCounter {
				http.Error(res, "Metrics has another type", http.StatusNotFound)
				return
			}
			valueString = strconv.FormatInt(typedValue, 10)
		}

		_, err := fmt.Fprintln(res, valueString)
		if err != nil {
			panic(err)
		}

		res.WriteHeader(http.StatusOK)
	}
}

func metricsRouter(storage MetricsStorage) chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.StripSlashes)
	// Need to route update requests to the same handler even if some named path components are absent
	// So we haven't found better way other than using such routing
	router.Route("/update", func(router chi.Router) {
		router.Post("/", handleUpdate(storage))
		router.Route("/{metricsType}", func(router chi.Router) {
			router.Post("/", handleUpdate(storage))
			router.Route("/{metricsName}", func(router chi.Router) {
				router.Post("/", handleUpdate(storage))
				router.Route("/{metricsValue}", func(router chi.Router) {
					router.Post("/", handleUpdate(storage))
				})
			})
		})
	})

	router.Get("/value/{metricsType}/{metricsName}", handleValue(storage))

	return router
}

func main() {
	err := http.ListenAndServe(`:8080`, metricsRouter(storage))
	if err != nil {
		panic(err)
	}
}
