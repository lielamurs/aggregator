package handlers

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lielamurs/aggregator/internal/config"
	"github.com/sirupsen/logrus"
)

func SetupRouter(handler *ApplicationHandler, cfg *config.Config, logger *logrus.Logger) *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	setupMiddleware(e, logger)
	setupRoutes(e, handler)

	return e
}

func setupMiddleware(e *echo.Echo, logger *logrus.Logger) {
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339}","id":"${id}","remote_ip":"${remote_ip}","host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}","status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}","bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
		Output: logger.Writer(),
	}))

	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 30 * time.Second,
	}))
}

func setupRoutes(e *echo.Echo, handler *ApplicationHandler) {
	e.GET("/health", handler.HealthCheck)

	v1 := e.Group("/api/v1")

	applications := v1.Group("/applications")
	applications.POST("", handler.SubmitApplication)
	applications.GET("/:id", handler.GetApplicationStatus)
}
