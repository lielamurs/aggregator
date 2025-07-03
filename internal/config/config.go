package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Banks    BanksConfig    `json:"banks"`
	Logging  LoggingConfig  `json:"logging"`
}

type ServerConfig struct {
	Port string `json:"port" env:"SERVER_PORT"`
	Host string `json:"host" env:"SERVER_HOST"`
}

type DatabaseConfig struct {
	Host         string `json:"host" env:"DB_HOST"`
	Port         string `json:"port" env:"DB_PORT"`
	User         string `json:"user" env:"DB_USER"`
	Password     string `json:"password" env:"DB_PASSWORD"`
	Name         string `json:"name" env:"DB_NAME"`
	SSLMode      string `json:"ssl_mode" env:"DB_SSLMODE"`
	MaxIdleConns int    `json:"max_idle_conns" env:"DB_MAX_IDLE_CONNS"`
	MaxOpenConns int    `json:"max_open_conns" env:"DB_MAX_OPEN_CONNS"`
	MaxLifetime  int    `json:"max_lifetime" env:"DB_MAX_LIFETIME"`
}

type BanksConfig struct {
	FastBank  FastBankConfig  `json:"fastbank"`
	SolidBank SolidBankConfig `json:"solidbank"`
}

type FastBankConfig struct {
	BaseURL      string `json:"base_url" env:"FASTBANK_BASE_URL"`
	Timeout      int    `json:"timeout" env:"FASTBANK_TIMEOUT"`
	MaxAttempts  int    `json:"max_attempts" env:"FASTBANK_MAX_ATTEMPTS"`
	PollInterval int    `json:"poll_interval" env:"FASTBANK_POLL_INTERVAL"`
}

type SolidBankConfig struct {
	BaseURL      string `json:"base_url" env:"SOLIDBANK_BASE_URL"`
	Timeout      int    `json:"timeout" env:"SOLIDBANK_TIMEOUT"`
	MaxAttempts  int    `json:"max_attempts" env:"SOLIDBANK_MAX_ATTEMPTS"`
	PollInterval int    `json:"poll_interval" env:"SOLIDBANK_POLL_INTERVAL"`
}

type LoggingConfig struct {
	Level  string `json:"level" env:"LOG_LEVEL"`
	Format string `json:"format" env:"LOG_FORMAT"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		logrus.Debug("No .env file found, using environment variables")
	}

	config := &Config{
		Server: ServerConfig{
			Port: getEnvOrDefault("SERVER_PORT", "8080"),
			Host: getEnvOrDefault("SERVER_HOST", "localhost"),
		},
		Database: DatabaseConfig{
			Host:         getEnvOrDefault("DB_HOST", "localhost"),
			Port:         getEnvOrDefault("DB_PORT", "5432"),
			User:         getEnvOrDefault("DB_USER", "postgres"),
			Password:     getEnvOrDefault("DB_PASSWORD", ""),
			Name:         getEnvOrDefault("DB_NAME", "aggregator"),
			SSLMode:      getEnvOrDefault("DB_SSLMODE", "disable"),
			MaxIdleConns: getEnvIntOrDefault("DB_MAX_IDLE_CONNS", 10),
			MaxOpenConns: getEnvIntOrDefault("DB_MAX_OPEN_CONNS", 100),
			MaxLifetime:  getEnvIntOrDefault("DB_MAX_LIFETIME", 3600),
		},
		Banks: BanksConfig{
			FastBank: FastBankConfig{
				BaseURL:      getEnvOrDefault("FASTBANK_BASE_URL", ""),
				Timeout:      getEnvIntOrDefault("FASTBANK_TIMEOUT", 30),
				MaxAttempts:  getEnvIntOrDefault("FASTBANK_MAX_ATTEMPTS", 30),
				PollInterval: getEnvIntOrDefault("FASTBANK_POLL_INTERVAL", 2),
			},
			SolidBank: SolidBankConfig{
				BaseURL:      getEnvOrDefault("SOLIDBANK_BASE_URL", ""),
				Timeout:      getEnvIntOrDefault("SOLIDBANK_TIMEOUT", 30),
				MaxAttempts:  getEnvIntOrDefault("SOLIDBANK_MAX_ATTEMPTS", 30),
				PollInterval: getEnvIntOrDefault("SOLIDBANK_POLL_INTERVAL", 2),
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
