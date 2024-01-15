package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Configuration struct {
	ServiceName    string
	LogLevel       string
	ApiURL         string
	SecureApiToken string
	ApiMaxRetries  int
}

var Config *Configuration

func LoadConfig() []error {

	var errs []error

	Config = &Configuration{
		ServiceName:    getEnv("SERVICE_NAME", "managed-clusters-onboard-tracking"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		ApiURL:         getEnv("API_URL", "https://secure.sysdig.com"),
		SecureApiToken: getEnv("SECURE_API_TOKEN", ""),
		ApiMaxRetries:  getIntEnv("API_MAX_RETRIES", 5),
	}

	if Config.ServiceName == "" {
		errs = append(errs, fmt.Errorf("SERVICE_NAME is missing"))
	}

	return errs
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getIntEnv(key string, defaultVal int) int {
	value, exists := os.LookupEnv(key)
	if exists {
		value, err := strconv.Atoi(value)
		if err != nil {
			log.Fatalf("Error converting env variable %s to integer. Please provide a integer convertable type.")
		}
		return value
	}
	return defaultVal
}
