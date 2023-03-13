package handler

import (
	"net/http"

	"kata-peya/internal/http/response"
	"kata-peya/internal/pet"
	"kata-peya/internal/tracer"

	"github.com/labstack/echo/v4"
)

type Pets struct {
	uc pet.UseCase
}

func NewPetsHandler(uc pet.UseCase) *Pets {
	return &Pets{uc: uc}
}

// GetAll godoc
// @Summary      List pets
// @Description  get pets
// @Tags         pets
// @Produce      json
// @Param        vaccinated    query     string  false  "get only vaccinated pets"  Format(boolean)
// @Success      200  {array}   response.Pet
// @Failure      400  {string}  string
// @Failure      404  {string}  string
// @Failure      500  {string}  string
// @Router       /pets [get]
func (p *Pets) GetAll(c echo.Context) error {
	ctx, span := tracer.Start(c.Request().Context(), "pets.handler.get_all")
	defer span.End()

	pets, err := p.uc.FindAll(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, response.Pets(pets))
}
