package configs

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// getEnv retrieves an environment variable or returns the fallback value.
func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && strings.TrimSpace(val) != "" {
		return strings.TrimSpace(val)
	}
	return fallback
}

// getEnvAsInt retrieves an environment variable as an int or returns the fallback value.
func getEnvAsInt(key string, fallback int) int {
	valStr := getEnv(key, "")
	if valStr == "" {
		return fallback
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return fallback
	}
	return val
}

// getEnvAsInt32 retrieves an environment variable as an int32 or returns the fallback value.
func getEnvAsInt32(key string, fallback int32) int32 {
	return int32(getEnvAsInt(key, int(fallback)))
}

// getEnvAsBool retrieves an environment variable as a boolean or returns the fallback value.
func getEnvAsBool(key string, fallback bool) bool {
	valStr := getEnv(key, "")
	if valStr == "" {
		return fallback
	}
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return fallback
	}
	return val
}

// getEnvAsDuration retrieves an environment variable as a time.Duration or returns the fallback value.
// It supports values like "15s", "10m", "24h", or plain numeric seconds.
func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	valStr := getEnv(key, "")
	if valStr == "" {
		return fallback
	}
	// Try parsing standard duration format (e.g., "10s", "1m", "24h")
	d, err := time.ParseDuration(valStr)
	if err == nil {
		return d
	}
	// Fallback to numeric seconds if plain number given
	if sec, err := strconv.Atoi(valStr); err == nil {
		return time.Duration(sec) * time.Second
	}
	return fallback
}

// getEnvAsSlice retrieves a comma-separated environment variable as a string slice.
func getEnvAsSlice(key string, fallback []string, sep string) []string {
	valStr := getEnv(key, "")
	if valStr == "" {
		return fallback
	}
	items := strings.Split(valStr, sep)
	var result []string
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	if len(result) == 0 {
		return fallback
	}
	return result
}
