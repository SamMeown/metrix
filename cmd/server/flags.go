package main

import "flag"

var flagAddress string

func parseFlags() {
	flag.StringVar(&flagAddress, "a", ":8080", "server address and port")
	flag.Parse()
}
