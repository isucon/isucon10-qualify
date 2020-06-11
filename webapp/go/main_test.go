package main_test

import (
	"fmt"
	main "github.com/isucon/isucon10-qualify/webapp/go"
	"os"
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

var MySQLConnectionData = main.NewMySQLConnectionEnv()
