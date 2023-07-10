package config_utils

import (
	"os"
	"strconv"
)

var LookupEnvString = os.LookupEnv

func LookupEnvInt(name string) (int, bool) {
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

func LookupEnvBool(name string) (bool, bool) {
	value, ok := os.LookupEnv(name)
	if !ok {
		return false, false
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		panic(err)
	}

	return boolValue, true
}
