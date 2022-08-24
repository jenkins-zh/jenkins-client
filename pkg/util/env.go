package util

import "os"

// GetEnvOrDefault returns a env or default value
func GetEnvOrDefault(key, defaultVal string) (result string) {
	result = defaultVal
	if val, ok := os.LookupEnv(key); ok {
		result = val
	}
	return
}
