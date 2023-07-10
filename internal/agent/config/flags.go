package config

import (
	"flag"
	"github.com/SamMeown/metrix/internal/utils/config_utils"
)

type Config struct {
	ServerBaseAddress string
	PollInterval      int
	ReportInterval    int
}

func Parse() Config {
	var config Config

	flag.StringVar(&config.ServerBaseAddress, "a", "localhost:8080", "metrics server address and port")
	flag.IntVar(&config.PollInterval, "p", 2, "metrics poll interval")
	flag.IntVar(&config.ReportInterval, "r", 10, "metrics report interval")

	flag.Parse()

	if address, ok := config_utils.LookupEnvString("ADDRESS"); ok {
		config.ServerBaseAddress = address
	}

	if pollInterval, ok := config_utils.LookupEnvInt("POLL_INTERVAL"); ok {
		config.PollInterval = pollInterval
	}

	if reportInterval, ok := config_utils.LookupEnvInt("REPORT_INTERVAL"); ok {
		config.ReportInterval = reportInterval
	}

	return config
}
