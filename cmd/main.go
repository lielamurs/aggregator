package main

import (
	"github.com/lielamurs/aggregator/internal/config"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load configuration")
	}

	// Setup logging
	setupLogging(cfg.Logging)

	logrus.WithFields(logrus.Fields{
		"server_host":        cfg.Server.Host,
		"server_port":        cfg.Server.Port,
		"fastbank_base_url":  cfg.Banks.FastBank.BaseURL,
		"solidbank_base_url": cfg.Banks.SolidBank.BaseURL,
	}).Info("Financing Application Aggregator initialized")

	logrus.Info("Configuration loaded successfully")
}

func setupLogging(cfg config.LoggingConfig) {
	// Set log level
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		logrus.WithError(err).Warn("Invalid log level, using info")
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	// Set log format
	switch cfg.Format {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	default:
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
}
