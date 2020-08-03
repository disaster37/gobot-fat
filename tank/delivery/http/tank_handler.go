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
	e.GET("/tanks", handler.GetState)

}

// GetState return the current tank state
func (h *TankHandler) GetState(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	level, err := h.dUsecase.Level(ctx)
	if err != nil {
		log.Errorf("Error when get Tank level: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, models.NewJSONAPIerror(
			"500",
			"Error when get Tank level",
			err.Error(),
			nil,
		))
	}
	volume, err := h.dUsecase.Volume(ctx)
	if err != nil {
		log.Errorf("Error when get Tank volume: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, models.NewJSONAPIerror(
			"500",
			"Error when get Tank level",
			err.Error(),
			nil,
		))
	}

	state := map[string]int{
		"level":  level,
		"volume": volume,
	}

	return c.JSON(http.StatusOK, models.JSONAPI{
		Data: models.JSONAPIData{
			Type:       "tanks",
			Id:         "state",
			Attributes: state,
		},
	})
}
