package env

import (
	"fmt"
	"os"
	"strconv"
)

func GetDBString(key, fallback string) string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" || port == "" || user == "" || pass == "" || dbname == "" {
		return fallback
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, dbname)
}

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
