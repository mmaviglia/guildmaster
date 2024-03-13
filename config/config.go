package config

import "os"

var DB_DSN = getEnv("DB_DSN", "")
var DISCORD_TOKEN = getEnv("DISCORD_TOKEN", "")
var ENV = getEnv("GO_ENV", "development")

// Get the environment variable associated with the key. If it
// does not exist, return value.
func getEnv(key, value string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return value
	}

	return v
}
