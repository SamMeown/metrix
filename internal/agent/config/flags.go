package config

import (
	"flag"
	"os"
	"strconv"
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

	if address, ok := lookupEnvString("ADDRESS"); ok {
		config.ServerBaseAddress = address
	}

	if pollInterval, ok := lookupEnvInt("POLL_INTERVAL"); ok {
		config.PollInterval = pollInterval
	}

	if reportInterval, ok := lookupEnvInt("REPORT_INTERVAL"); ok {
		config.ReportInterval = reportInterval
	}

	return config
}

var lookupEnvString = os.LookupEnv

func lookupEnvInt(name string) (int, bool) {
	value, ok := os.LookupEnv(name)
	if !ok {
		return 0, false
	}

	intValue, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		panic(err)
	}

	return int(intValue), true
}
