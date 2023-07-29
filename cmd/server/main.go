package main

import (
	"context"
	"database/sql"
	"github.com/SamMeown/metrix/internal/storage"
	"github.com/SamMeown/metrix/internal/storage/pg"
	"github.com/SamMeown/metrix/internal/storage/retryable"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/SamMeown/metrix/internal/server"
	"github.com/SamMeown/metrix/internal/server/config"
	"github.com/SamMeown/metrix/internal/server/saver"
)

func main() {
	ctx := context.Background()
	serverConfig := config.Parse()

	var metricsStorage storage.MetricsStorage
	var storageSaver *saver.MetricsStorageSaver
	if len(serverConfig.DatabaseDSN) > 0 {
		db, err := sql.Open("pgx", serverConfig.DatabaseDSN)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		pgStorage := pg.NewStorage(db)
		err = pgStorage.Bootstrap(ctx)
		if err != nil {
			panic(err)
		}

		metricsStorage = retryable.NewStorage(pgStorage, pg.IsRetryableError)
	} else {
		metricsStorage = storage.New()

		var err error
		storageSaver, err = saver.NewMetricsStorageSaver(metricsStorage, serverConfig.StoragePath)
		if err != nil {
			panic(err)
		}
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

	server.Run(ctx, serverConfig, metricsStorage, storageSaver)
}
