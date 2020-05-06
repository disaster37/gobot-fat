package http

import (
	"context"
	"net/http"

	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/tfp"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

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
		return c.JSON(http.StatusInternalServerError, models.NewJSONAPIerror(
			"500",
			"Error when get TFP state",
			err.Error(),
			nil,
		))
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
		return c.JSON(500, models.NewJSONAPIerror(
			"500",
			"Error when start pond pump",
			err.Error(),
			nil,
		))
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
		log.Errorf("Error when post start_pond_pump: %s", err.Error())
		return c.JSON(500, models.NewJSONAPIerror(
			"500",
			"Error when start pond pump",
			err.Error(),
			nil,
		))
	}

	err = h.dUsecase.UVC1(ctx, true)
	if err != nil {
		log.Errorf("Error when post start_uvc1: %s", err.Error())
		return c.JSON(500, models.NewJSONAPIerror(
			"500",
			"Error when start uvc1",
			err.Error(),
			nil,
		))
	}

	err = h.dUsecase.UVC2(ctx, true)
	if err != nil {
		log.Errorf("Error when post start_uvc2: %s", err.Error())
		return c.JSON(500, models.NewJSONAPIerror(
			"500",
			"Error when start uvc2",
			err.Error(),
			nil,
		))
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
		return c.JSON(500, models.NewJSONAPIerror(
			"500",
			"Error when stop pond pomp",
			err.Error(),
			nil,
		))
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
		return c.JSON(500, models.NewJSONAPIerror(
			"500",
			"Error when start waterfall pump",
			err.Error(),
			nil,
		))
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
		return c.JSON(500, models.NewJSONAPIerror(
			"500",
			"Error when stop waterfall pump",
			err.Error(),
			nil,
		))
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
		return c.JSON(500, models.NewJSONAPIerror(
			"500",
			"Error when start uvc1",
			err.Error(),
			nil,
		))
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
		return c.JSON(500, models.NewJSONAPIerror(
			"500",
			"Error when stop uvc1",
			err.Error(),
			nil,
		))
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
		return c.JSON(500, models.NewJSONAPIerror(
			"500",
			"Error when start uvc2",
			err.Error(),
			nil,
		))
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
		return c.JSON(500, models.NewJSONAPIerror(
			"500",
			"Error when stop uvc2",
			err.Error(),
			nil,
		))
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
		return c.JSON(500, models.NewJSONAPIerror(
			"500",
			"Error when start pond bubble",
			err.Error(),
			nil,
		))
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
		return c.JSON(500, models.NewJSONAPIerror(
			"500",
			"Error when stop pond bubble",
			err.Error(),
			nil,
		))
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
		return c.JSON(500, models.NewJSONAPIerror(
			"500",
			"Error when start filter bubble",
			err.Error(),
			nil,
		))
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
		return c.JSON(500, models.NewJSONAPIerror(
			"500",
			"Error when stop filter bubble",
			err.Error(),
			nil,
		))
	}

	return c.NoContent(http.StatusOK)
}
