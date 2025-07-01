package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Server  ServerConfig  `json:"server"`
	Banks   BanksConfig   `json:"banks"`
	Logging LoggingConfig `json:"logging"`
}

type ServerConfig struct {
	Port string `json:"port" env:"SERVER_PORT"`
	Host string `json:"host" env:"SERVER_HOST"`
}

type BanksConfig struct {
	FastBank  FastBankConfig  `json:"fastbank"`
	SolidBank SolidBankConfig `json:"solidbank"`
}

type FastBankConfig struct {
	BaseURL string `json:"base_url" env:"FASTBANK_BASE_URL"`
	Timeout int    `json:"timeout" env:"FASTBANK_TIMEOUT"`
}

type SolidBankConfig struct {
	BaseURL string `json:"base_url" env:"SOLIDBANK_BASE_URL"`
	Timeout int    `json:"timeout" env:"SOLIDBANK_TIMEOUT"`
}

type LoggingConfig struct {
	Level  string `json:"level" env:"LOG_LEVEL"`
	Format string `json:"format" env:"LOG_FORMAT"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Try to load .env file if it exists
	if err := godotenv.Load(); err != nil {
		logrus.Debug("No .env file found, using environment variables")
	}

	config := &Config{
		Server: ServerConfig{
			Port: getEnvOrDefault("SERVER_PORT", "8080"),
			Host: getEnvOrDefault("SERVER_HOST", "localhost"),
		},
		Banks: BanksConfig{
			FastBank: FastBankConfig{
				BaseURL: getEnvOrDefault("FASTBANK_BASE_URL", ""),
				Timeout: getEnvIntOrDefault("FASTBANK_TIMEOUT", 30),
			},
			SolidBank: SolidBankConfig{
				BaseURL: getEnvOrDefault("SOLIDBANK_BASE_URL", ""),
				Timeout: getEnvIntOrDefault("SOLIDBANK_TIMEOUT", 30),
			},
		},
		Logging: LoggingConfig{
			Level:  getEnvOrDefault("LOG_LEVEL", "info"),
			Format: getEnvOrDefault("LOG_FORMAT", "json"),
		},
	}

	return config, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
