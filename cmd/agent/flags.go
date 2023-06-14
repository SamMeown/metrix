package main

import "flag"

var serverBaseAddress string
var pollInterval int
var reportInterval int

func parseFlags() {
	flag.StringVar(&serverBaseAddress, "a", "localhost:8080", "metrics server address and port")
	flag.IntVar(&pollInterval, "p", 2, "metrics poll interval")
	flag.IntVar(&reportInterval, "r", 10, "metrics report interval")

	flag.Parse()
}
