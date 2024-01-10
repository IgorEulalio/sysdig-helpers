package config

import (
	"fmt"
	"os"
)

type Configuration struct {
	ServiceName    string
	LogLevel       string
	ApiURL         string
	SecureApiToken string
}

var Config *Configuration

func LoadConfig() []error {

	var errs []error

	Config = &Configuration{
		ServiceName:    getEnv("SERVICE_NAME", ""),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		ApiURL:         getEnv("API_URL", "https://secure.sysdig.com"),
		SecureApiToken: getEnv("SECURE_API_TOKEN", ""),
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
