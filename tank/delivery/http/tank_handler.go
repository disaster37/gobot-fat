package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/tank"
	"github.com/google/jsonapi"
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
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	values, err := h.dUsecase.Tanks(ctx)
	if err != nil {
		log.Errorf("Error when list tanks values: %s", err.Error())
		c.Response().WriteHeader(http.StatusInternalServerError)
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when list tanks value",
				Detail: err.Error(),
			},
		})
	}
	tanks := make([]*models.Tank, 0, len(values))
	for _, tank := range values {
		tanks = append(tanks, tank)
	}

	c.Response().WriteHeader(http.StatusOK)
	return jsonapi.MarshalPayload(c.Response(), tanks)
}

// GetTankValues return the tank value
func (h *TankHandler) GetTankValues(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	name := c.Param("id")
	log.Debugf("Name: %s", name)

	value, err := h.dUsecase.Tank(ctx, name)
	if err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when get tank value",
				Detail: err.Error(),
			},
		})
	}

	c.Response().WriteHeader(http.StatusOK)
	return jsonapi.MarshalOnePayloadEmbedded(c.Response(), value)
}
