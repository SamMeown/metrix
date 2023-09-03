package agent

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/SamMeown/metrix/internal/agent/client"
	"github.com/SamMeown/metrix/internal/agent/config"
	"github.com/SamMeown/metrix/internal/agent/metrics"
	"github.com/SamMeown/metrix/internal/logger"
)

var wg sync.WaitGroup
var m sync.Mutex
var done = make(chan struct{})

func startRefreshMetricsLoop(conf config.Config, collector metrics.Collector) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Duration(conf.PollInterval) * time.Second)
		for {
			m.Lock()
			logger.Log.Infoln("Refreshing metrics...")
			collector.CollectMetrics()
			m.Unlock()
			select {
			case <-ticker.C:
				//continue
			case <-done:
				return
			}
		}
	}()
}

func startReportMetricsLoop(conf config.Config, collector metrics.Collector, client *client.MetricsClient) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Duration(conf.ReportInterval) * time.Second)
		for {
			m.Lock()
			logger.Log.Infoln("Reporting metrics...")
			client.ReportAllMetrics(collector.GetMetrics())
			collector.ResetCounters()
			m.Unlock()
			select {
			case <-ticker.C:
				//continue
			case <-done:
				return
			}
		}
	}()
}

func monitorDoneSignals() {
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	<-shutdownCh
	close(done)
}

func Start(conf config.Config, collector metrics.Collector, client *client.MetricsClient) {
	err := logger.Initialize("debug")
	if err != nil {
		panic(err)
	}
	defer logger.Log.Sync()

	startRefreshMetricsLoop(conf, collector)
	startReportMetricsLoop(conf, collector, client)

	go monitorDoneSignals()

	wg.Wait()

	logger.Log.Infoln("Agent stopped")
}
