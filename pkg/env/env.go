// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Really simple env variable helper, with default values
// ----------------------------------------------------------------------------

package env

import (
	"os"
	"strconv"
)

// Internal function to fetch environmental variable or return default
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// GetEnvString is a simple helper function to read an environment variable
// as a string or return a default value.
func GetEnvString(key string, defaultVal string) string {
	return getEnv(key, defaultVal)
}

// GetEnvInt is a simple helper function to read an environment variable
// as an integer or return a default value.
func GetEnvInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// GetEnvFloat is a simple helper function to read an environment variable
// as a float or return a default value.
func GetEnvFloat(key string, defaultVal float64) float64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return value
	}

	return defaultVal
}

// GetEnvBool is a simple helper function to read an environment variable
// as a boolean or return a default value.
func GetEnvBool(key string, defaultVal bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}

	return defaultVal
}
