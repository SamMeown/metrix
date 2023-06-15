package main

import (
	"flag"
	"os"
)

var flagAddress string

func parseFlags() {
	flag.StringVar(&flagAddress, "a", ":8080", "server address and port")
	flag.Parse()

	if envAddress, ok := os.LookupEnv("ADDRESS"); ok {
		flagAddress = envAddress
	}
}
