package agent

import (
	"fmt"
	"github.com/SamMeown/metrix/internal/agent/client"
	"github.com/SamMeown/metrix/internal/agent/config"
	"github.com/SamMeown/metrix/internal/agent/metrics"
	"time"
)

func Start(conf config.Config, collector metrics.Collector, client *client.MetricsClient) {
	var secondsElapsed int64

	for {
		if secondsElapsed%int64(conf.PollInterval) == 0 {
			fmt.Println("Refreshing metrics...")
			collector.CollectMetrics()
		}

		if secondsElapsed%int64(conf.ReportInterval) == 0 {
			fmt.Println("Reporting metrics...")
			client.ReportAllMetrics(collector.Collection())
			collector.ResetCollectsCount()
		}

		secondsElapsed++
		time.Sleep(1 * time.Second)
	}
}