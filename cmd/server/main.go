package main

import (
	"github.com/SamMeown/metrix/internal/server"
	"github.com/SamMeown/metrix/internal/server/config"
	"github.com/SamMeown/metrix/internal/server/saver"
	"github.com/SamMeown/metrix/internal/storage"
)

func main() {
	serverConfig := config.Parse()
	metricsStorage := storage.New()
	storageSaver, err := saver.NewMetricsStorageSaver(metricsStorage, serverConfig.StoragePath)
	if err != nil {
		panic(err)
	}
	server.Start(serverConfig, metricsStorage, storageSaver)
}
