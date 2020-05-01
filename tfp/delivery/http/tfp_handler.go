package http

import (
	"context"
	"net/http"

	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/tfp"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"error"`
	Code    int    `json:"error_code"`
}

// TFPHandler  represent the httphandler for tfp
type TFPHandler struct {
	dUsecase tfp.Usecase
}

// NewTFPHandler will initialize the TFP_config/ resources endpoint
func NewTFPHandler(e *echo.Group, us tfp.Usecase) {
	handler := &TFPHandler{
		dUsecase: us,
	}
	e.POST("/tfps/action/start_pond_pump", handler.StartPondPump)
	e.POST("/tfps/action/start_pond_pump_with_uvc", handler.StartPondPumpWithUVC)
	e.POST("/tfps/action/stop_pond_pump", handler.StopPondPump)
	e.POST("/tfps/action/start_waterfall_pump", handler.StartWaterfallPump)
	e.POST("/tfps/action/stop_waterfall_pump", handler.StopWaterfallPump)
	e.POST("/tfps/action/start_uvc1", handler.StartUVC1)
	e.POST("/tfps/action/stop_uvc1", handler.StopUVC1)
	e.POST("/tfps/action/start_uvc2", handler.StartUVC2)
	e.POST("/tfps/action/stop_uvc2", handler.StopUVC2)
	e.POST("/tfps/action/start_pond_bubble", handler.StartPondBubble)
	e.POST("/tfps/action/stop_pond_bubble", handler.StopPondBubble)
	e.POST("/tfps/action/start_filter_bubble", handler.StartFilterBubble)
	e.POST("/tfps/action/stop_filter_bubble", handler.StopFilterBubble)
	e.GET("/tfps", handler.GetState)

}

// GetState return the current state of TFP
func (h TFPHandler) GetState(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	state, err := h.dUsecase.GetState(ctx)

	if err != nil {
		log.Errorf("Error when get TFP state: %s", err.Error())
		return c.JSON(500, models.JSONAPI{
			Errors: []models.JSONAPIError{
				models.JSONAPIError{
					Status: "500",
					Title:  "Error when get TFP state",
					Detail: err.Error(),
				},
			},
		})
	}

	return c.JSON(http.StatusOK, models.JSONAPI{
		Data: models.JSONAPIData{
			Type:       "tfps",
			Id:         "state",
			Attributes: state,
		},
	})
}

// StartPondPump start pond pump
func (h TFPHandler) StartPondPump(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := h.dUsecase.PondPump(ctx, true)

	if err != nil {
		log.Errorf("Error when post start_pond_pump: %s", err.Error())
		return c.JSON(500, ResponseError{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

// StartPondPumpWithUVC start pond pump and then UVC
func (h TFPHandler) StartPondPumpWithUVC(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := h.dUsecase.PondPump(ctx, true)

	if err != nil {
		log.Errorf("Error when post start_pond_pump_with_uvc: %s", err.Error())
		return c.JSON(500, ResponseError{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	err = h.dUsecase.UVC1(ctx, true)
	if err != nil {
		log.Errorf("Error when post start_pond_pump_with_uvc: %s", err.Error())
		return c.JSON(500, ResponseError{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	err = h.dUsecase.UVC2(ctx, true)
	if err != nil {
		log.Errorf("Error when post start_pond_pump_with_uvc: %s", err.Error())
		return c.JSON(500, ResponseError{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

// StopPondPump stop pond pump
func (h TFPHandler) StopPondPump(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := h.dUsecase.PondPump(ctx, false)

	if err != nil {
		log.Errorf("Error when post stop_pond_pump: %s", err.Error())
		return c.JSON(500, ResponseError{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

// StartWaterfallPump start waterfall pump
func (h TFPHandler) StartWaterfallPump(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := h.dUsecase.WaterfallPump(ctx, true)

	if err != nil {
		log.Errorf("Error when post start_waterfall_pump: %s", err.Error())
		return c.JSON(500, ResponseError{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

// StopWaterfallPump stop waterfall pump
func (h TFPHandler) StopWaterfallPump(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := h.dUsecase.WaterfallPump(ctx, false)

	if err != nil {
		log.Errorf("Error when post stop_waterfall_pump: %s", err.Error())
		return c.JSON(500, ResponseError{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

// StartUVC1 start UVC1
func (h TFPHandler) StartUVC1(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := h.dUsecase.UVC1(ctx, true)

	if err != nil {
		log.Errorf("Error when post start_uvc1: %s", err.Error())
		return c.JSON(500, ResponseError{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

// StopUVC1 stop UVC1
func (h TFPHandler) StopUVC1(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := h.dUsecase.UVC1(ctx, false)

	if err != nil {
		log.Errorf("Error when post stop_uvc1: %s", err.Error())
		return c.JSON(500, ResponseError{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

// StartUVC2 start UVC2
func (h TFPHandler) StartUVC2(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := h.dUsecase.UVC2(ctx, true)

	if err != nil {
		log.Errorf("Error when post start_uvc2: %s", err.Error())
		return c.JSON(500, ResponseError{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

// StopUVC2 stop UVC2
func (h TFPHandler) StopUVC2(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := h.dUsecase.UVC2(ctx, false)

	if err != nil {
		log.Errorf("Error when post stop_uvc2: %s", err.Error())
		return c.JSON(500, ResponseError{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

// StartPondBubble start pond bubble
func (h TFPHandler) StartPondBubble(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := h.dUsecase.PondBubble(ctx, true)

	if err != nil {
		log.Errorf("Error when post start_pond_bubble: %s", err.Error())
		return c.JSON(500, ResponseError{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

// StopPondBubble stop pond bubble
func (h TFPHandler) StopPondBubble(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := h.dUsecase.PondBubble(ctx, false)

	if err != nil {
		log.Errorf("Error when post stop_pond_bubble: %s", err.Error())
		return c.JSON(500, ResponseError{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

// StartFilterBubble start filter bubble
func (h TFPHandler) StartFilterBubble(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := h.dUsecase.FilterBubble(ctx, true)

	if err != nil {
		log.Errorf("Error when post start_filter_bubble: %s", err.Error())
		return c.JSON(500, ResponseError{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

// StopFilterBubble stop filter bubble
func (h TFPHandler) StopFilterBubble(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := h.dUsecase.FilterBubble(ctx, false)

	if err != nil {
		log.Errorf("Error when post stop_filter_bubble: %s", err.Error())
		return c.JSON(500, ResponseError{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}
