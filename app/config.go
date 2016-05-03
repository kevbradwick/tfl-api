package app

import (
	"fmt"
	"os"
)

//
// Get an environment variable and if it doesn't exist, return a default value
//
func getEnv(name, defaultVal string) string {
	if os.Getenv(name) != "" {
		return os.Getenv(name)
	}
	return defaultVal
}

//
// The host+port of the mongod instance to connect to
//
func MongoUrl() string {
	host := getEnv("mongohost", "localhost")
	port := getEnv("mongoport", "27017")
	return fmt.Sprintf("%s:%s", host, port)
}
