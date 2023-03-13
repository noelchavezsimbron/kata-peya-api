package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthCheck godoc
// @Summary Show the status of server.
// @Description Get the status of server.
// @Tags healthCheck
// @Accept */*
// @Produce json
// @Success 200 {object} string
// @Router /health [get]
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
