package http

import (
	"context"
	"net/http"

	"github.com/disaster37/gobot-fat/dfp"
	"github.com/disaster37/gobot-fat/models"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// DFPHandler  represent the httphandler for dfp
type DFPHandler struct {
	dUsecase dfp.Usecase
}

// NewDFPHandler will initialize the DFP/ resources endpoint
func NewDFPHandler(e *echo.Group, us dfp.Usecase) {
	handler := &DFPHandler{
		dUsecase: us,
	}
	e.POST("/dfps/action/start", handler.Start)
	e.POST("/dfps/action/stop", handler.Stop)
	e.POST("/dfps/action/wash", handler.Wash)
	e.POST("/dfps/action/manual_start_drum", handler.ManualStartDrum)
	e.POST("/dfps/action/manual_stop_drum", handler.ManualStopDrum)
	e.POST("/dfps/action/manual_start_pump", handler.ManualStartPump)
	e.POST("/dfps/action/manual_stop_pump", handler.ManualStopPump)
	e.POST("/dfps/action/set_security", handler.SetSecurity)
	e.POST("/dfps/action/unset_security", handler.UnsetSecurity)
	e.POST("/dfps/action/set_emergency_stop", handler.SetEmergencyStop)
	e.POST("/dfps/action/unset_emergency_stop", handler.UnsetEmergencyStop)
	e.GET("/dfps", handler.GetState)

}

// GetState return the current state of TFP
func (h DFPHandler) GetState(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	state, err := h.dUsecase.GetState(ctx)

	if err != nil {
		log.Errorf("Error when get DFP state: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, models.NewJSONAPIerror(
			"500",
			"Error when get DFP state",
			err.Error(),
			nil,
		))
	}

	return c.JSON(http.StatusOK, models.JSONAPI{
		Data: models.JSONAPIData{
			Type:       "dfps",
			Id:         "state",
			Attributes: state,
		},
	})
}

// Start put DFP on auto mode
func (h DFPHandler) Start(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := h.dUsecase.Start(ctx)

	if err != nil {
		log.Errorf("Error when post start: %s", err.Error())
		return c.JSON(500, models.NewJSONAPIerror(
			"500",
			"Error when start DFP",
			err.Error(),
			nil,
		))
	}

	return c.NoContent(http.StatusNoContent)
}

// Stop put stop mode on DFP
func (h DFPHandler) Stop(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := h.dUsecase.Stop(ctx)

	if err != nil {
		log.Errorf("Error when post stop: %s", err.Error())
		return c.JSON(500, models.NewJSONAPIerror(
			"500",
			"Error when stop DFP",
			err.Error(),
			nil,
		))
	}

	return c.NoContent(http.StatusNoContent)
}

// Wash force wash cycle
func (h DFPHandler) Wash(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := h.dUsecase.Wash(ctx)

	if err != nil {
		log.Errorf("Error when post wash: %s", err.Error())
		return c.JSON(500, models.NewJSONAPIerror(
			"500",
			"Error when force wash",
			err.Error(),
			nil,
		))
	}

	return c.NoContent(http.StatusNoContent)
}

// ManualStartDrum start drum motor
func (h DFPHandler) ManualStartDrum(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := h.dUsecase.ManualDrum(ctx, true)

	if err != nil {
		log.Errorf("Error when post manual_start_drum: %s", err.Error())
		return c.JSON(500, models.NewJSONAPIerror(
			"500",
			"Error when start drum motor",
			err.Error(),
			nil,
		))
	}

	return c.NoContent(http.StatusNoContent)
}

// ManualStopDrum stop drum motor
func (h DFPHandler) ManualStopDrum(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := h.dUsecase.ManualDrum(ctx, false)

	if err != nil {
		log.Errorf("Error when post manual_stop_drum: %s", err.Error())
		return c.JSON(500, models.NewJSONAPIerror(
			"500",
			"Error when stop drum motor",
			err.Error(),
			nil,
		))
	}

	return c.NoContent(http.StatusNoContent)
}

// ManualStartPump start pump
func (h DFPHandler) ManualStartPump(c echo.Context) error {
	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}

	err := h.dUsecase.ManualPump(ctx, true)

	if err != nil {
		log.Errorf("Error when post manual_start_pump: %s", err.Error())
		return c.JSON(500, models.NewJSONAPIerror(
			"500",
			"Error when start pump",
			err.Error(),
			nil,
		))
	}

	return c.NoContent(http.StatusNoContent)
}

// ManualStopPump stop pump
func (h DFPHandler) ManualStopPump(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := h.dUsecase.ManualPump(ctx, false)

	if err != nil {
		log.Errorf("Error when post manual_stop_pump: %s", err.Error())
		return c.JSON(500, models.NewJSONAPIerror(
			"500",
			"Error when stop pump",
			err.Error(),
			nil,
		))
	}

	return c.NoContent(http.StatusNoContent)
}

// SetSecurity permit to set security
func (h DFPHandler) SetSecurity(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := h.dUsecase.Security(ctx, true)

	if err != nil {
		log.Errorf("Error when post set_security: %s", err.Error())
		return c.JSON(500, models.NewJSONAPIerror(
			"500",
			"Error when set security",
			err.Error(),
			nil,
		))
	}

	return c.NoContent(http.StatusNoContent)
}

// UnsetSecurity permit to remove security
func (h DFPHandler) UnsetSecurity(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := h.dUsecase.Security(ctx, false)

	if err != nil {
		log.Errorf("Error when post unset_security: %s", err.Error())
		return c.JSON(500, models.NewJSONAPIerror(
			"500",
			"Error when unset security",
			err.Error(),
			nil,
		))
	}

	return c.NoContent(http.StatusNoContent)
}

// SetEmergencyStop permit to set emergency
func (h DFPHandler) SetEmergencyStop(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := h.dUsecase.EmergencyStop(ctx, true)

	if err != nil {
		log.Errorf("Error when post set_emergency_stop: %s", err.Error())
		return c.JSON(500, models.NewJSONAPIerror(
			"500",
			"Error when set emergency stop",
			err.Error(),
			nil,
		))
	}

	return c.NoContent(http.StatusNoContent)
}

// UnsetEmergencyStop permit to set emergency
func (h DFPHandler) UnsetEmergencyStop(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := h.dUsecase.EmergencyStop(ctx, false)

	if err != nil {
		log.Errorf("Error when post unset_emergency_stop: %s", err.Error())
		return c.JSON(500, models.NewJSONAPIerror(
			"500",
			"Error when unset emergency stop",
			err.Error(),
			nil,
		))
	}

	return c.NoContent(http.StatusNoContent)
}
