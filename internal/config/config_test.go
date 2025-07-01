package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	config, err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.Server.Port != "8080" {
		t.Errorf("Expected default port 8080, got %s", config.Server.Port)
	}

	if config.Server.Host != "localhost" {
		t.Errorf("Expected default host localhost, got %s", config.Server.Host)
	}

	if config.Logging.Level != "info" {
		t.Errorf("Expected default log level info, got %s", config.Logging.Level)
	}

	if config.Logging.Format != "json" {
		t.Errorf("Expected default log format json, got %s", config.Logging.Format)
	}

	if config.Banks.FastBank.Timeout != 30 {
		t.Errorf("Expected default FastBank timeout 30, got %d", config.Banks.FastBank.Timeout)
	}

	if config.Banks.SolidBank.Timeout != 30 {
		t.Errorf("Expected default SolidBank timeout 30, got %d", config.Banks.SolidBank.Timeout)
	}
}

func TestLoadWithEnvironmentVariables(t *testing.T) {
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("SERVER_HOST", "0.0.0.0")
	os.Setenv("FASTBANK_BASE_URL", "https://fastbank.example.com")
	os.Setenv("SOLIDBANK_BASE_URL", "https://solidbank.example.com")
	os.Setenv("FASTBANK_TIMEOUT", "60")
	os.Setenv("SOLIDBANK_TIMEOUT", "45")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("LOG_FORMAT", "text")

	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("FASTBANK_BASE_URL")
		os.Unsetenv("SOLIDBANK_BASE_URL")
		os.Unsetenv("FASTBANK_TIMEOUT")
		os.Unsetenv("SOLIDBANK_TIMEOUT")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("LOG_FORMAT")
	}()

	config, err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.Server.Port != "9090" {
		t.Errorf("Expected port 9090, got %s", config.Server.Port)
	}

	if config.Server.Host != "0.0.0.0" {
		t.Errorf("Expected host 0.0.0.0, got %s", config.Server.Host)
	}

	if config.Banks.FastBank.BaseURL != "https://fastbank.example.com" {
		t.Errorf("Expected FastBank URL https://fastbank.example.com, got %s", config.Banks.FastBank.BaseURL)
	}

	if config.Banks.SolidBank.BaseURL != "https://solidbank.example.com" {
		t.Errorf("Expected SolidBank URL https://solidbank.example.com, got %s", config.Banks.SolidBank.BaseURL)
	}

	if config.Banks.FastBank.Timeout != 60 {
		t.Errorf("Expected FastBank timeout 60, got %d", config.Banks.FastBank.Timeout)
	}

	if config.Banks.SolidBank.Timeout != 45 {
		t.Errorf("Expected SolidBank timeout 45, got %d", config.Banks.SolidBank.Timeout)
	}

	if config.Logging.Level != "debug" {
		t.Errorf("Expected log level debug, got %s", config.Logging.Level)
	}

	if config.Logging.Format != "text" {
		t.Errorf("Expected log format text, got %s", config.Logging.Format)
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	result := getEnvOrDefault("TEST_VAR", "default_value")
	if result != "test_value" {
		t.Errorf("Expected test_value, got %s", result)
	}

	result = getEnvOrDefault("NON_EXISTING_VAR", "default_value")
	if result != "default_value" {
		t.Errorf("Expected default_value, got %s", result)
	}

	os.Setenv("EMPTY_VAR", "")
	defer os.Unsetenv("EMPTY_VAR")

	result = getEnvOrDefault("EMPTY_VAR", "default_value")
	if result != "default_value" {
		t.Errorf("Expected default_value for empty env var, got %s", result)
	}
}

func TestGetEnvIntOrDefault(t *testing.T) {
	os.Setenv("TEST_INT_VAR", "42")
	defer os.Unsetenv("TEST_INT_VAR")

	result := getEnvIntOrDefault("TEST_INT_VAR", 10)
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}

	result = getEnvIntOrDefault("NON_EXISTING_INT_VAR", 10)
	if result != 10 {
		t.Errorf("Expected 10, got %d", result)
	}

	os.Setenv("INVALID_INT_VAR", "not_a_number")
	defer os.Unsetenv("INVALID_INT_VAR")

	result = getEnvIntOrDefault("INVALID_INT_VAR", 10)
	if result != 10 {
		t.Errorf("Expected 10 for invalid int, got %d", result)
	}

	os.Setenv("EMPTY_INT_VAR", "")
	defer os.Unsetenv("EMPTY_INT_VAR")

	result = getEnvIntOrDefault("EMPTY_INT_VAR", 10)
	if result != 10 {
		t.Errorf("Expected 10 for empty int env var, got %d", result)
	}

	os.Setenv("ZERO_INT_VAR", "0")
	defer os.Unsetenv("ZERO_INT_VAR")

	result = getEnvIntOrDefault("ZERO_INT_VAR", 10)
	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}
