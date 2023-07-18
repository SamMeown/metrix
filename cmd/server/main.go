package main

import (
	"database/sql"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/SamMeown/metrix/internal/server"
	"github.com/SamMeown/metrix/internal/server/config"
	"github.com/SamMeown/metrix/internal/server/saver"
	"github.com/SamMeown/metrix/internal/storage"
)

func main() {
	serverConfig := config.Parse()

	db, err := sql.Open("pgx", serverConfig.DatabaseDSN)
	if err != nil {
		panic(err)
	}
	defer db.Close()

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

	server.Run(serverConfig, metricsStorage, storageSaver, db)
}
