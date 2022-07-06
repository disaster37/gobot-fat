package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/disaster37/gobot-fat/dfp"
	"github.com/google/jsonapi"
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
	e.POST("/dfps/action/set_disable_security", handler.SetDisableSecurity)
	e.POST("/dfps/action/unset_disable_security", handler.UnsetDisableSecurity)
	e.POST("/dfps/action/set_emergency_stop", handler.SetEmergencyStop)
	e.POST("/dfps/action/unset_emergency_stop", handler.UnsetEmergencyStop)
	e.GET("/dfps", handler.GetState)
	e.GET("/dfps/io", handler.GetIO)

}

// GetState return the current state of DFP
func (h DFPHandler) GetState(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	state, err := h.dUsecase.GetState(ctx)
	if err != nil {
		log.Errorf("Error when get DFP state: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when get DFP state",
				Detail: err.Error(),
			},
		})
	}

	log.Infof("DFPState: %s", spew.Sdump(&state))

	c.Response().WriteHeader(http.StatusOK)
	return jsonapi.MarshalOnePayloadEmbedded(c.Response(), &state)
}

// GetIO return the current IO of DFP
func (h DFPHandler) GetIO(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	io, err := h.dUsecase.GetIO(ctx)
	if err != nil {
		log.Errorf("Error when get DFP IO: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when get DFP IO",
				Detail: err.Error(),
			},
		})
	}

	log.Infof("DFPIO: %s", spew.Sdump(&io))

	c.Response().WriteHeader(http.StatusOK)
	return jsonapi.MarshalOnePayloadEmbedded(c.Response(), &io)
}

// Start put DFP on auto mode
func (h DFPHandler) Start(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.Start(ctx)

	if err != nil {
		log.Errorf("Error when post start: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when start DFP",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// Stop put stop mode on DFP
func (h DFPHandler) Stop(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.Stop(ctx)

	if err != nil {
		log.Errorf("Error when post stop: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when stop DFP",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// Wash force wash cycle
func (h DFPHandler) Wash(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.Wash(ctx)

	if err != nil {
		log.Errorf("Error when post wash: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when force wash",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// ManualStartDrum start drum motor
func (h DFPHandler) ManualStartDrum(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.ManualDrum(ctx, true)

	if err != nil {
		log.Errorf("Error when post manual_start_drum: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when start drum motor",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// ManualStopDrum stop drum motor
func (h DFPHandler) ManualStopDrum(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.ManualDrum(ctx, false)

	if err != nil {
		log.Errorf("Error when post manual_stop_drum: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when stop drum motor",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// ManualStartPump start pump
func (h DFPHandler) ManualStartPump(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.ManualPump(ctx, true)

	if err != nil {
		log.Errorf("Error when post manual_start_pump: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when start pump",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// ManualStopPump stop pump
func (h DFPHandler) ManualStopPump(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.ManualPump(ctx, false)

	if err != nil {
		log.Errorf("Error when post manual_stop_pump: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when stop pump",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// SetSecurity permit to set security
func (h DFPHandler) SetSecurity(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.Security(ctx, true)

	if err != nil {
		log.Errorf("Error when post set_security: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when set security",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// UnsetSecurity permit to remove security
func (h DFPHandler) UnsetSecurity(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.Security(ctx, false)

	if err != nil {
		log.Errorf("Error when post unset_security: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when unset security",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// SetEmergencyStop permit to set emergency
func (h DFPHandler) SetEmergencyStop(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.EmergencyStop(ctx, true)

	if err != nil {
		log.Errorf("Error when post set_emergency_stop: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when set emergency stop",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// UnsetEmergencyStop permit to set emergency
func (h DFPHandler) UnsetEmergencyStop(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.EmergencyStop(ctx, false)

	if err != nil {
		log.Errorf("Error when post unset_emergency_stop: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when unset emergency stop",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// SetDisableSecurity permit to disable security
func (h DFPHandler) SetDisableSecurity(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.DisableSecurity(ctx, true)

	if err != nil {
		log.Errorf("Error when post set_disable_security: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when disable security",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// UnsetDisableSecurity permit to remove security
func (h DFPHandler) UnsetDisableSecurity(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.DisableSecurity(ctx, false)

	if err != nil {
		log.Errorf("Error when post unset_disable_security: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when enable security",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}
