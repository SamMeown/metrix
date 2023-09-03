package main

import (
	"github.com/SamMeown/metrix/internal/agent"
	"github.com/SamMeown/metrix/internal/agent/client"
	"github.com/SamMeown/metrix/internal/agent/config"
	"github.com/SamMeown/metrix/internal/agent/metrics"
	"github.com/SamMeown/metrix/internal/crypto/signer"
	"github.com/SamMeown/metrix/internal/storage"
)

func main() {
	agentConfig := config.Parse()
	mStorage := storage.NewMemStorage()
	mCollector := metrics.NewCollector(mStorage)
	mSigner := signer.New(agentConfig.SignKey)
	mClient := client.NewMetricsClient(agentConfig.ServerBaseAddress, mSigner)

	agent.Start(agentConfig, mCollector, mClient)
}
