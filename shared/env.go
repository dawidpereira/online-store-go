package shared

import (
	"os"
	"strconv"
	"time"
)

func GetString(key, defaultValue string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return val
}

func GetInt(key string, defaultValue int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}

	return intVal
}

func GetBool(key string, defaultValue bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return defaultValue
	}

	return boolVal
}

func GetDuration(key string, defaultValue time.Duration) time.Duration {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	durationVal, err := time.ParseDuration(val)
	if err != nil {
		return defaultValue
	}

	return durationVal
}
