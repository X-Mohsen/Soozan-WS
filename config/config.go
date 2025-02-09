package config

import "os"

var ListenPort = getEnv("WS_SERVER_PORT", "8080")

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
