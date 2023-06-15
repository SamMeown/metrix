package main

import (
	"flag"
	"os"
	"strconv"
)

var serverBaseAddress string
var pollInterval int
var reportInterval int

func parseFlags() {
	flag.StringVar(&serverBaseAddress, "a", "localhost:8080", "metrics server address and port")
	flag.IntVar(&pollInterval, "p", 2, "metrics poll interval")
	flag.IntVar(&reportInterval, "r", 10, "metrics report interval")

	flag.Parse()

	tryEnvString("ADDRESS", func(value string) {
		serverBaseAddress = value
	})

	tryEnvInt("POLL_INTERVAL", func(value int) {
		pollInterval = value
	})

	tryEnvInt("REPORT_INTERVAL", func(value int) {
		reportInterval = value
	})
}

func tryEnv(name string, callback func(value string)) {
	if value, ok := os.LookupEnv(name); ok {
		callback(value)
	}
}

var tryEnvString = tryEnv

func tryEnvInt(name string, callback func(value int)) {
	tryEnv(name, func(strValue string) {
		value, err := strconv.ParseInt(strValue, 10, 32)
		if err != nil {
			panic(err)
		}

		callback(int(value))
	})
}
