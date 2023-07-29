package agent

import (
	"time"

	"github.com/SamMeown/metrix/internal/agent/client"
	"github.com/SamMeown/metrix/internal/agent/config"
	"github.com/SamMeown/metrix/internal/agent/metrics"
	"github.com/SamMeown/metrix/internal/logger"
)

func Start(conf config.Config, collector metrics.Collector, client *client.MetricsClient) {
	err := logger.Initialize("info")
	if err != nil {
		panic(err)
	}
	defer logger.Log.Sync()

	var secondsElapsed int64

	for {
		if secondsElapsed%int64(conf.PollInterval) == 0 {
			logger.Log.Infoln("Refreshing metrics...")
			collector.CollectMetrics()
		}

		if secondsElapsed%int64(conf.ReportInterval) == 0 {
			logger.Log.Infoln("Reporting metrics...")
			client.ReportAllMetrics(collector.Collection())
			collector.ResetCollectsCount()
		}

		secondsElapsed++
		time.Sleep(1 * time.Second)
	}
}
