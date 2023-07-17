package main

import (
	"os"
	"os/signal"
	"syscall"

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

	go func() {
		gracefulShutdown := make(chan os.Signal, 1)
		signal.Notify(gracefulShutdown,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)

		<-gracefulShutdown
		server.Stop()
	}()

	server.Run(serverConfig, metricsStorage, storageSaver)
}
