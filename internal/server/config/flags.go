package config

import (
	"flag"
	"os"
)

type Config struct {
	Address string
}

func Parse() (config Config) {
	flag.StringVar(&config.Address, "a", ":8080", "server address and port")
	flag.Parse()

	if envAddress, ok := os.LookupEnv("ADDRESS"); ok {
		config.Address = envAddress
	}

	return
}
