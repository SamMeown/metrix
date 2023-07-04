package main

import (
	"github.com/SamMeown/metrix/internal/agent"
	"github.com/SamMeown/metrix/internal/agent/client"
	"github.com/SamMeown/metrix/internal/agent/config"
	"github.com/SamMeown/metrix/internal/agent/metrics"
	"github.com/SamMeown/metrix/internal/storage"
)

func main() {
	agentConfig := config.Parse()
	mStorage := storage.New()
	mCollector := metrics.NewCollector(mStorage)
	mClient := client.NewMetricsClient(agentConfig.ServerBaseAddress)

	agent.Start(agentConfig, mCollector, mClient)
}
