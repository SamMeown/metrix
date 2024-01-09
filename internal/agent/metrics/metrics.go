package metrics

import (
	"context"
	"github.com/SamMeown/metrix/internal/logger"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"math/rand"
	"reflect"
	"runtime"
	"sync"

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
	GetMetrics() storage.MetricsStorageItems
	ResetCounters()
}

func NewCollector(storage *storage.MemStorage) Collector {
	return &metricsCollector{
		memStatsMetricsNames: memStatsMetricsNames[:],
		collection:           storage,
	}
}

type metricsCollector struct {
	m                    sync.Mutex
	wg                   sync.WaitGroup
	memStatsMetricsNames []string
	collection           *storage.MemStorage
}

func (mg *metricsCollector) GetMetrics() storage.MetricsStorageItems {
	mg.m.Lock()
	defer mg.m.Unlock()
	allMetrics, _ := mg.collection.GetAll(context.Background())
	return allMetrics
}

func (mg *metricsCollector) CollectMetrics() {
	mg.wg.Add(2)

	go func() {
		defer mg.wg.Done()
		mg.collectMemstatMetrics()
		mg.collectAdditionalMetrics()
	}()

	go func() {
		defer mg.wg.Done()
		mg.collectPsutilMetrics()
	}()

	mg.wg.Wait()
}

func (mg *metricsCollector) collectAdditionalMetrics() {
	ctx := context.Background()

	mg.m.Lock()
	defer mg.m.Unlock()
	mg.collection.SetGauge(ctx, "RandomValue", float64(rand.Int()))
	mg.collection.SetCounter(ctx, "PollCount", 1)
}

func (mg *metricsCollector) collectPsutilMetrics() {
	ctx := context.Background()
	mg.collectPsutilMemMetrics(ctx)
	mg.collectPsutilCPUMetrics(ctx)
}

func (mg *metricsCollector) collectPsutilMemMetrics(ctx context.Context) {
	memStat, err := mem.VirtualMemory()
	if err != nil {
		logger.Log.Errorf("Failed to get mem info: %s", err)
		return
	}

	mg.m.Lock()
	defer mg.m.Unlock()
	mg.collection.SetGauge(ctx, "TotalMemory", gauge(memStat.Total))
	mg.collection.SetGauge(ctx, "FreeMemory", gauge(memStat.Available))
}

func (mg *metricsCollector) collectPsutilCPUMetrics(ctx context.Context) {
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		logger.Log.Errorf("Failed to get cpu percentage: %s", err)
		return
	}

	mg.m.Lock()
	defer mg.m.Unlock()
	mg.collection.SetGauge(ctx, "CPUutilization1", cpuPercent[0])
}

func (mg *metricsCollector) collectMemstatMetrics() {
	var memstats runtime.MemStats
	runtime.ReadMemStats(&memstats)

	ctx := context.Background()

	mg.m.Lock()
	defer mg.m.Unlock()
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

func (mg *metricsCollector) ResetCounters() {
	mg.m.Lock()
	defer mg.m.Unlock()
	ctx := context.Background()
	mg.collection.ResetCounters(ctx)
}
