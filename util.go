package main

import "os"

func getEnv(key, defaultValue string) string {
	val, isSet := os.LookupEnv(key)
	if isSet {
		return val
	}
	return defaultValue
}
