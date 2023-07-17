package config

import (
	"flag"

	"github.com/SamMeown/metrix/internal/utils/config_utils"
)

type Config struct {
	Address       string
	StoreInterval int
	StoragePath   string
	Restore       bool
}

func Parse() (config Config) {
	flag.StringVar(&config.Address, "a", ":8080", "server address and port")
	flag.IntVar(&config.StoreInterval, "i", 300, "metrics saving time interval")
	flag.StringVar(&config.StoragePath, "f", "/tmp/metrics-db.json", "storage dump file path")
	flag.BoolVar(&config.Restore, "r", true, "should restore from saved dump on start")
	flag.Parse()

	if envAddress, ok := configutils.LookupEnvString("ADDRESS"); ok {
		config.Address = envAddress
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

	return
}
