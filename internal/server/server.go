package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/SamMeown/metrix/internal/logger"
	"github.com/SamMeown/metrix/internal/models"
	"github.com/SamMeown/metrix/internal/server/config"
	middlewares "github.com/SamMeown/metrix/internal/server/middleware"
	"github.com/SamMeown/metrix/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"strconv"
)

var tableTemplate = `
<!DOCTYPE html>
<html>
<head>
<style>
table {
  font-family: arial, sans-serif;
  border-collapse: collapse;
  width: 100%%;
}

td, th {
  border: 1px solid #dddddd;
  text-align: left;
  padding: 8px;
}

tr:nth-child(even) {
  background-color: #dddddd;
}
</style>
</head>
<body>

<h2>Metrics</h2>

<table>
  <tr>
    <th>Name</th>
    <th>Value</th>
  </tr>
  %s
</table>

</body>
</html>
`

var tableRowTemlate = `<tr>
    <td>%s</td>
    <td>%v</td>
  </tr>`

func handleUpdateJSON(mStorage storage.MetricsStorage) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		var metrics models.Metrics
		var buf bytes.Buffer

		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(buf.Bytes(), &metrics)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		logger.Log.Debugf("Body: %+v", metrics)

		if metrics.ID == "" {
			http.Error(res, "No metrics name", http.StatusNotFound)
			return
		}

		switch metrics.MType {
		case storage.MetricsTypeGauge:
			if metrics.Value != nil {
				mStorage.SetGauge(metrics.ID, *metrics.Value)
			} else {
				http.Error(res, "No metrics value", http.StatusBadRequest)
				return
			}
		case storage.MetricsTypeCounter:
			if metrics.Delta != nil {
				mStorage.SetCounter(metrics.ID, *metrics.Delta)
			} else {
				http.Error(res, "No metrics value", http.StatusBadRequest)
				return
			}
		default:
			http.Error(res, "Wrong metrics type", http.StatusBadRequest)
			return
		}

		response := metrics
		response.Delta = nil
		switch response.MType {
		case storage.MetricsTypeGauge:
			value, _ := mStorage.GetGauge(response.ID)
			response.Value = value
		case storage.MetricsTypeCounter:
			counter, _ := mStorage.GetCounter(response.ID)
			value := float64(*counter)
			response.Value = &value
		}

		resp, err := json.Marshal(response)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}

		res.WriteHeader(http.StatusOK)
		_, err = res.Write(resp)
		if err != nil {
			logger.Log.Errorf("Failed to write response body")
		}
	}
}

func handleUpdate(mStorage storage.MetricsStorage) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/plain; charset=utf-8")

		var metricsType = chi.URLParam(req, "metricsType")
		var metricsName = chi.URLParam(req, "metricsName")
		var metricsValueStr = chi.URLParam(req, "metricsValue")

		logger.Log.Debugf("metricsType: %s, metricsName: %s, metricsValue: %s\n", metricsType, metricsName, metricsValueStr)

		if metricsType == "" {
			http.Error(res, "Wrong number of data components", http.StatusBadRequest)
			return
		}

		if metricsType != storage.MetricsTypeGauge &&
			metricsType != storage.MetricsTypeCounter {
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

		switch metricsType {
		case storage.MetricsTypeGauge:
			if metricsValue, convErr := strconv.ParseFloat(metricsValueStr, 64); convErr == nil {
				mStorage.SetGauge(metricsName, metricsValue)
			} else {
				http.Error(res, "Can not parse metrics value", http.StatusBadRequest)
				return
			}
		case storage.MetricsTypeCounter:
			if metricsValue, convErr := strconv.ParseInt(metricsValueStr, 10, 64); convErr == nil {
				mStorage.SetCounter(metricsName, metricsValue)
			} else {
				http.Error(res, "Can not parse metrics value", http.StatusBadRequest)
				return
			}
		}

		res.WriteHeader(http.StatusOK)
	}
}

func handleValueJSON(mStorage storage.MetricsStorage) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		var request models.Metrics
		var buf bytes.Buffer

		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(buf.Bytes(), &request)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		logger.Log.Debugf("Body: %+v", request)

		if request.ID == "" {
			http.Error(res, "No metrics name", http.StatusBadRequest)
			return
		}

		response := request

		switch request.MType {
		default:
			http.Error(res, "Wrong metrics type", http.StatusBadRequest)
			return
		case storage.MetricsTypeGauge:
			value, _ := mStorage.GetGauge(request.ID)
			if value == nil {
				http.Error(res, "Metrics not found", http.StatusNotFound)
				return
			}
			response.Value = value
		case storage.MetricsTypeCounter:
			value, _ := mStorage.GetCounter(request.ID)
			if value == nil {
				http.Error(res, "Metrics not found", http.StatusNotFound)
				return
			}
			response.Delta = value
		}

		resp, err := json.Marshal(response)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusOK)
		_, err = res.Write(resp)
		if err != nil {
			logger.Log.Errorf("Failed to write response body")
		}
	}
}

func handleValue(mStorage storage.MetricsStorage) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/plain; charset=utf-8")

		var metricsType = chi.URLParam(req, "metricsType")
		var metricsName = chi.URLParam(req, "metricsName")

		var valueString string
		switch metricsType {
		default:
			http.Error(res, "Wrong metrics type", http.StatusBadRequest)
			return
		case storage.MetricsTypeGauge:
			value, _ := mStorage.GetGauge(metricsName)
			if value == nil {
				http.Error(res, "Metrics not found", http.StatusNotFound)
				return
			}
			valueString = strconv.FormatFloat(*value, 'f', -1, 64)
		case storage.MetricsTypeCounter:
			value, _ := mStorage.GetCounter(metricsName)
			if value == nil {
				http.Error(res, "Metrics not found", http.StatusNotFound)
				return
			}
			valueString = strconv.FormatInt(*value, 10)
		}

		res.WriteHeader(http.StatusOK)

		_, err := fmt.Fprintln(res, valueString)
		if err != nil {
			panic(err)
		}
	}
}

func handleRoot(mStorage storage.MetricsStorage) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/html; charset=utf-8")

		var rows string
		snapshot, _ := mStorage.GetAll()
		for name, value := range snapshot.Gauges {
			rows += fmt.Sprintf(tableRowTemlate, name, value)
		}
		for name, value := range snapshot.Counters {
			rows += fmt.Sprintf(tableRowTemlate, name, value)
		}

		table := fmt.Sprintf(tableTemplate, rows)

		res.WriteHeader(http.StatusOK)

		_, err := fmt.Fprintln(res, table)
		if err != nil {
			panic(err)
		}
	}
}

func metricsRouter(mStorage storage.MetricsStorage) chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.StripSlashes, middlewares.Logging)
	// Need to route update requests to the same handler even if some named path components are absent
	// So we haven't found better way other than using such routing
	router.Route("/update", func(router chi.Router) {
		router.Post("/", handleUpdateJSON(mStorage))
		router.Route("/{metricsType}", func(router chi.Router) {
			router.Post("/", handleUpdate(mStorage))
			router.Route("/{metricsName}", func(router chi.Router) {
				router.Post("/", handleUpdate(mStorage))
				router.Route("/{metricsValue}", func(router chi.Router) {
					router.Post("/", handleUpdate(mStorage))
				})
			})
		})
	})

	router.Get("/value/{metricsType}/{metricsName}", handleValue(mStorage))

	router.Post("/value", handleValueJSON(mStorage))

	router.Get("/", handleRoot(mStorage))

	return router
}

func Start(conf config.Config, mStorage storage.MetricsStorage) {
	err := logger.Initialize("debug")
	if err != nil {
		panic(err)
	}
	defer logger.Log.Sync()

	err = http.ListenAndServe(conf.Address, metricsRouter(mStorage))
	if err != nil {
		panic(err)
	}
}
