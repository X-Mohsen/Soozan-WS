package config

import "os"

var APIBASEURL = getEnv("API_BASE_URL", "http://localhost:8000")
var PublicKeyPath = getEnv("PUBLIC_KEY_PATH", "public_key.pem")
var ListenPort = getEnv("WS_SERVER_PORT", "8080")

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
