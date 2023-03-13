package http

import (
	"kata-peya/internal/http/handler"

	"github.com/labstack/echo/v4"
)

type Api struct {
	petsHandler *handler.Pets
}

func NewApi(petsHandler *handler.Pets) *Api {
	return &Api{petsHandler: petsHandler}
}

func (api *Api) Routes(r *echo.Group) {
	r.GET("/health", handler.HealthCheck)
	r.GET("/pets", api.petsHandler.GetAll)
}
