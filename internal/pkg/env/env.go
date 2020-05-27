package env

import (
	"os"
	"strconv"
)

// GetEnvInt get integer from environment variable or default value
func GetEnvInt(key string, defaultValue int) int {
	v, err := strconv.ParseInt(os.Getenv(key), 10, 32)
	if err != nil {
		return int(v)
	}
	return defaultValue
}
