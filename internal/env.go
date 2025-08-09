package env

import (
	"os"
	"strconv"
)

func GetString(key, fallback string) string {

	val, ok := os.LookupEnv(key)

	if !ok {
		return fallback
	}

	return val
}

func GetInt(key string, fallback int) int {

	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	val1, _ := strconv.Atoi(val)

	return val1

}
