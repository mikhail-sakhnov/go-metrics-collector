package agent

import (
	"log"
	"os"
	"strconv"
)

func mustNotErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func getEnvInt(key string, defaultValue int) int {
	v, err := strconv.ParseInt(os.Getenv(key), 10, 32)
	if err != nil {
		return int(v)
	}
	return defaultValue
}
