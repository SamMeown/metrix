package metrics

import (
	"context"
	"math/rand"
	"reflect"
	"runtime"

	"github.com/SamMeown/metrix/internal/storage"
)

type gauge = float64
type counter = int64

var memStatsMetricsNames = []string{
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

type Collector interface {
	CollectMetrics()
	Collection() storage.MetricsStorageGetter
	ResetCollectsCount()
}

func NewCollector(storage storage.MetricsStorage) Collector {
	return &metricsCollector{
		memStatsMetricsNames: memStatsMetricsNames[:],
		collection:           storage,
	}
}

type metricsCollector struct {
	collectsCount        int64
	memStatsMetricsNames []string
	collection           storage.MetricsStorage
}

func (mg *metricsCollector) Collection() storage.MetricsStorageGetter {
	return mg.collection
}

func (mg *metricsCollector) CollectMetrics() {
	mg.collectsCount++

	ctx := context.Background()
	mg.collectMemstatMetrics(ctx)

	mg.collection.SetGauge(ctx, "RandomValue", float64(rand.Int()))
	mg.collection.SetCounter(ctx, "PollCount", mg.collectsCount)
}

func (mg *metricsCollector) collectMemstatMetrics(ctx context.Context) {
	var memstats runtime.MemStats
	runtime.ReadMemStats(&memstats)

	for _, name := range mg.memStatsMetricsNames {
		memstatsValue := reflect.ValueOf(memstats)
		field := memstatsValue.FieldByName(name)

		var fieldValue gauge
		if field.CanFloat() {
			fieldValue = gauge(field.Float())
		} else if field.CanUint() {
			fieldValue = gauge(field.Uint())
		} else {
			fieldValue = gauge(field.Int())
		}

		mg.collection.SetGauge(ctx, name, fieldValue)
	}
}

func (mg *metricsCollector) ResetCollectsCount() {
	mg.collectsCount = 0
}
