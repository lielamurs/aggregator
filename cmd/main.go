package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lielamurs/aggregator/internal/config"
	"github.com/lielamurs/aggregator/internal/handlers"
	"github.com/lielamurs/aggregator/internal/repository"
	"github.com/lielamurs/aggregator/internal/services"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load configuration")
	}

	// Setup logging
	logger := setupLogging(cfg.Logging)

	logger.WithFields(logrus.Fields{
		"server_host": cfg.Server.Host,
		"server_port": cfg.Server.Port,
		"db_host":     cfg.Database.Host,
		"db_port":     cfg.Database.Port,
		"db_name":     cfg.Database.Name,
	}).Info("Financing Application Aggregator initializing")

	// Initialize database connection
	db, err := repository.NewConnection(cfg.Database, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize database connection")
	}
	logger.Info("Database connection established")

	// Initialize repositories
	applicationsRepo := repository.NewApplicationsRepository(db.DB)
	offersRepo := repository.NewOffersRepository(db.DB)
	bankSubmissionsRepo := repository.NewBankSubmissionsRepository(db.DB)
	logger.Info("Repositories initialized")

	// Initialize bank services
	var bankServices []services.BankService

	if cfg.Banks.FastBank.BaseURL != "" {
		fastBankService := services.NewFastBankService(cfg.Banks.FastBank, logger)
		bankServices = append(bankServices, fastBankService)
		logger.Info("FastBank service initialized")
	} else {
		logger.Warn("FastBank BaseURL not configured, skipping FastBank integration")
	}

	if cfg.Banks.SolidBank.BaseURL != "" {
		solidBankService := services.NewSolidBankService(cfg.Banks.SolidBank, logger)
		bankServices = append(bankServices, solidBankService)
		logger.Info("SolidBank service initialized")
	} else {
		logger.Warn("SolidBank BaseURL not configured, skipping SolidBank integration")
	}

	if len(bankServices) == 0 {
		logger.Fatal("No bank services configured. Please configure at least one bank service.")
	}

	// Initialize application service with repositories
	applicationService := services.NewApplicationService(
		applicationsRepo,
		offersRepo,
		bankSubmissionsRepo,
		bankServices,
		logger,
	)
	logger.Info("Application service initialized")

	// Initialize handlers
	applicationHandler := handlers.NewApplicationHandler(applicationService, logger)
	logger.Info("HTTP handlers initialized")

	// Setup router
	router := handlers.SetupRouter(applicationHandler, cfg, logger)
	logger.Info("HTTP router configured")

	// Start server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.WithFields(logrus.Fields{
			"host": cfg.Server.Host,
			"port": cfg.Server.Port,
		}).Info("Starting HTTP server")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Failed to start HTTP server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Give the server 30 seconds to shutdown gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("Server forced to shutdown")
	} else {
		logger.Info("Server shutdown completed")
	}

	// Close database connection
	if err := db.Close(); err != nil {
		logger.WithError(err).Error("Failed to close database connection")
	} else {
		logger.Info("Database connection closed")
	}
}

func setupLogging(cfg config.LoggingConfig) *logrus.Logger {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		logger.WithError(err).Warn("Invalid log level, using info")
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

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

	return logger
}
