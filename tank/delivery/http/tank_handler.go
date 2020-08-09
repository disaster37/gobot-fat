package http

import (
	"context"
	"net/http"

	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/tank"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// TankHandler represent the httphandler for tank
type TankHandler struct {
	dUsecase tank.Usecase
}

// NewTFPHandler will initialize the TFP_config/ resources endpoint
func NewTankHandler(e *echo.Group, us tank.Usecase) {
	handler := &TankHandler{
		dUsecase: us,
	}
	e.GET("/tanks", handler.GetTanksValues)
	e.GET("/tanks/:id", handler.GetTankValues)

}

// GetTanksValues return the current tanks state
func (h *TankHandler) GetTanksValues(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	values, err := h.dUsecase.Tanks(ctx)
	if err != nil {
		log.Errorf("Error when get Tanks values: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, models.NewJSONAPIerror(
			"500",
			"Error when get Tanks values",
			err.Error(),
			nil,
		))
	}

	data := make([]models.JSONAPIData, 0, 1)
	for name, value := range values {
		data = append(data, models.JSONAPIData{
			Type:       "tanks",
			Id:         name,
			Attributes: value,
		})
	}

	return c.JSON(http.StatusOK, models.JSONAPI{
		Data: data,
	})
}

// GetTankValues return the tank value
func (h *TankHandler) GetTankValues(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	name := c.Param("id")
	log.Debugf("Name: %s", name)

	value, err := h.dUsecase.Tank(ctx, name)
	if err != nil {
		log.Errorf("Error when get Tank values: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, models.NewJSONAPIerror(
			"500",
			"Error when get Tanks values",
			err.Error(),
			nil,
		))
	}

	return c.JSON(http.StatusOK, models.JSONAPI{
		Data: models.JSONAPIData{
			Type:       "tanks",
			Id:         name,
			Attributes: value,
		},
	})
}
