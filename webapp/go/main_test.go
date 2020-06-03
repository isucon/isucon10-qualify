package main_test

import (
    "os"
    "fmt"
)

func getEnv(key, defaultValue string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	return defaultValue
}

func getURL() string {
	port := getEnv("API_PORT", "1323")
	host := getEnv("API_HOST", "localhost")
	url := fmt.Sprintf("http://%s:%s", host, port)
	return url
}

