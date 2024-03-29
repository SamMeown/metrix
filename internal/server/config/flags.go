package config

import (
	"flag"

	"github.com/SamMeown/metrix/internal/utils/config_utils"
)

type Config struct {
	Address       string
	DatabaseDSN   string
	StoreInterval int
	StoragePath   string
	Restore       bool
	SignKey       string
}

func Parse() (config Config) {
	flag.StringVar(&config.Address, "a", ":8080", "server address and port")
	flag.StringVar(&config.DatabaseDSN, "d", "", "database dsn")
	flag.IntVar(&config.StoreInterval, "i", 300, "metrics saving time interval")
	flag.StringVar(&config.StoragePath, "f", "/tmp/metrics-db.json", "storage dump file path")
	flag.BoolVar(&config.Restore, "r", true, "should restore from saved dump on start")
	flag.StringVar(&config.SignKey, "k", "", "signature key")
	flag.Parse()

	if envAddress, ok := configutils.LookupEnvString("ADDRESS"); ok {
		config.Address = envAddress
	}

	if envDatabaseDSN, ok := configutils.LookupEnvString("DATABASE_DSN"); ok {
		config.DatabaseDSN = envDatabaseDSN
	}

	if envStoreInterval, ok := configutils.LookupEnvInt("STORE_INTERVAL"); ok {
		config.StoreInterval = envStoreInterval
	}

	if envStoragePath, ok := configutils.LookupEnvString("FILE_STORAGE_PATH"); ok {
		config.StoragePath = envStoragePath
	}

	if envRestore, ok := configutils.LookupEnvBool("RESTORE"); ok {
		config.Restore = envRestore
	}

	if envKey, ok := configutils.LookupEnvString("KEY"); ok {
		config.SignKey = envKey
	}

	return
}
