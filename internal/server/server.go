package server

import (
	"fmt"
	"github.com/SamMeown/metrix/internal/server/config"
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

func handleUpdate(mStorage storage.MetricsStorage) func(http.ResponseWriter, *http.Request) {
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

func handleValue(mStorage storage.MetricsStorage) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/plain; charset=utf-8")

		path := req.URL.Path
		fmt.Printf("Path: %s\n", path)

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

		_, err := fmt.Fprintln(res, valueString)
		if err != nil {
			panic(err)
		}

		res.WriteHeader(http.StatusOK)
	}
}

func handleRoot(mStorage storage.MetricsStorage) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		var rows string
		metrics, _ := mStorage.GetAll()
		for name, value := range metrics {
			rows += fmt.Sprintf(tableRowTemlate, name, value)
		}

		table := fmt.Sprintf(tableTemplate, rows)

		_, err := fmt.Fprintln(res, table)
		if err != nil {
			panic(err)
		}

		res.Header().Set("Content-Type", "text/html; charset=utf-8")
		res.WriteHeader(http.StatusOK)
	}
}

func metricsRouter(mStorage storage.MetricsStorage) chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.StripSlashes)
	// Need to route update requests to the same handler even if some named path components are absent
	// So we haven't found better way other than using such routing
	router.Route("/update", func(router chi.Router) {
		router.Post("/", handleUpdate(mStorage))
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

	router.Get("/", handleRoot(mStorage))

	return router
}

func Start(conf config.Config, mStorage storage.MetricsStorage) {
	err := http.ListenAndServe(conf.Address, metricsRouter(mStorage))
	if err != nil {
		panic(err)
	}
}
