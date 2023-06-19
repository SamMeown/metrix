package main

import (
	"github.com/SamMeown/metrix/internal/server"
	"github.com/SamMeown/metrix/internal/server/config"
	"github.com/SamMeown/metrix/internal/storage"
)

func main() {
	serverConfig := config.Parse()
	metricsStorage := storage.New()
	server.Start(serverConfig, metricsStorage)
}
